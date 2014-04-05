package gieba

import (
    "os"
    "math"
    "bufio"
    "strconv"
    "strings"
    "unicode/utf8"
)


type TreeItem struct {
    freq float64    // if zero, means from root to here there is no str exists
    next map[rune]*TreeItem
}


type Freq struct {
    data  *TreeItem
    total float64   // accumulation of all freq
    minf  float64   // minimun of all freq
    hmm   *HMM
}


func FreqNew() (f *Freq){
    f = &Freq{&TreeItem{0.0, make(map[rune]*TreeItem)},
              0.0, math.MaxFloat64, nil}
    return
}


type FI struct {
    // source string index src[i:j]
    // so the next FI will be src[j:k], i < j < k
    i, j int
    freq float64
    one bool    // one rune or not
}


type FIS struct {
    data []*FI
}


func (f *Freq)LoadFile(infn string) (err error){
    fin, err := os.Open(infn)
    if err != nil { return }
    defer fin.Close()

    scanner := bufio.NewScanner(fin)
    for scanner.Scan() {
        cols := strings.Fields(scanner.Text())
        o, e := strconv.ParseFloat(cols[1], 32)
        if e != nil { return e }
        p := f.data
        for _, c := range cols[0] {
            q, ok := p.next[c]
            if ! ok {
                q = &TreeItem{0.0, make(map[rune]*TreeItem)}
                p.next[c] = q
            }
            p = q
        }
        //if p.freq != 0.0 {
        //    log.Log("Duplicate entry", p.freq, o)
        //}
        p.freq = o
    }
    return
}


func (f *Freq)getDAGbyTree(src string) (dag []*FIS, idxs []int){
    dag = make([]*FIS, len(src))
    idxs = make([]int, 0, len(src) + 1)
    for i := range src {
        idxs = append(idxs, i)
        p := f.data
        l := FIS{make([]*FI, 0)}
        for j, c := range src[i:] {
            q, ok := p.next[c]
            //if ! ok { break }
            r := 0.0
            if ok {
                p = q
                r = p.freq
            }
            if r == 0 && j == 0 {
                r = f.minf  // always add first rune, even not found
            }
            if r > 0 {
                e := i + j + utf8.RuneLen(c)
                fi := FI{i, e, r, j == 0}
                l.data = append(l.data, &fi)
            }
            if ! ok { break }
        }
        dag[i] = &l
    }
    idxs = append(idxs, len(src))
    return
}


func (f *Freq)Cut(src string) (ret []string) {
    src_len := len(src)
    dag, _ := f.getDAGbyTree(src)
    cand := make([]*FI, src_len + 1)
    // mark ending
    cand[src_len] = &FI{src_len, src_len, 0.0, true}
    // calulate max freq for each i
    for i := len(dag) - 1; i  >= 0; i-- {
        if dag[i] == nil {continue}
        maxf := 0.0
        for _, fi := range dag[i].data {
            newf := fi.freq + cand[fi.j].freq
            if newf > maxf {
                cand[i] = fi
            }
        }
    }
    // build return
    ret = make([]string, 0, src_len / 2)
    oi := 0
    var d *FI
    for d = cand[0]; d.i != d.j; d = cand[d.j] {
        if f.hmm != nil {
            if d.one && d.j < src_len {continue}
            if oi < d.i {
                // d is more then one sigle rune
                // process buffered
                // TODO, check again if exists in Tree
                n := f.hmm.Cut(src[oi:d.i])
                ret = append(ret, n...)
            }
        }
        ret = append(ret, src[d.i:d.j])
        oi = d.j
    }
    return
}
