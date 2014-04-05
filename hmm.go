package gieba

// http://en.wikipedia.org/wiki/Viterbi_algorithm

import (
    "os"
    "log"
    "fmt"
    "path"
    "io/ioutil"
    "encoding/json"
)


type HMM struct {
    states  []string
    start   map[string]float64
    trans   map[string]map[string]float64
    emit    map[string]map[string]float64
    prev    map[string][]string
}


func loadOneJson(fn string, dst interface{}) {
    fin, err := os.Open(fn)
    if err != nil {
        log.Fatal(err)
    }
    defer fin.Close()
    dat, err := ioutil.ReadAll(fin)
    err = json.Unmarshal(dat, dst)
    if err != nil {
        log.Fatal(err)
    }
}


func HMMNew(fn string) (v *HMM) {
    v = &HMM{make([]string, 0, 5),
             make(map[string]float64),
             make(map[string]map[string]float64),
             make(map[string]map[string]float64),
             make(map[string][]string),
                }
    if fn != "" {
        v.LoadJsons(fn)
    }
    return
}


func (v *HMM)LoadJsons(datadir string) {
    loadOneJson(path.Join(datadir, "prob_start") + ".json", &(v.start))
    loadOneJson(path.Join(datadir, "prob_trans") + ".json", &(v.trans))
    loadOneJson(path.Join(datadir, "prob_emit") + ".json", &(v.emit))
    loadOneJson(path.Join(datadir, "prob_prev") + ".json", &(v.prev))
    for k, _ := range v.start {
        v.states = append(v.states, k)
    }
}


const (
    MIN_FLOAT float64 = -3.14e100
)


func get(m map[string]float64, k string, d float64) float64 {
    r, ok := m[k]
    if ok {return r}
    return d
}


func show(n string, v interface{}) {
    j, _ := json.Marshal(v)
    fmt.Println(n, string(j))
}


func (v *HMM)Viterbi(obs string) (float64, []string, []int) {
    // prepare index list to operate string as rune list
    idxs := make([]int, 0, len(obs) + 1)
    for idx, _ := range obs {idxs = append(idxs, idx)}
    // if only one rune, return it
    if len(idxs) == 1 {return 0.0, nil, idxs}
    idxs = append(idxs, len(obs))
    // init
    V := make([]map[string]float64, len(idxs))
    path := make(map[string][]string)
    V[0] = make(map[string]float64)
    oidx := idxs[1]
    for _, y := range v.states {
        V[0][y] = v.start[y] + get(v.emit[y], obs[:oidx], MIN_FLOAT)
        path[y] = []string{y}
    }
    t, t1 := 0, 1
    for _, idx := range idxs[2:] {
        newpath := make(map[string][]string)
        V[t1] = make(map[string]float64)
        for _, y := range v.states {
            var prob0 float64
            var state0 string
            em_p := get(v.emit[y], obs[oidx:idx], MIN_FLOAT)
            for _, y0 := range v.prev[y] {
                p := V[t][y0] + get(v.trans[y0], y, MIN_FLOAT) + em_p
                if p > prob0 || prob0 == 0 { prob0, state0 = p, y0 }
            }
            V[t1][y] = prob0
            newpath[y] = append(newpath[y], path[state0]...)
            newpath[y] = append(newpath[y], y)
        }
        t, t1, oidx, path = t1, t1 + 1, idx, newpath
    }
    var prob float64
    var state string
    for _, y := range []string{"E", "S"} {
        p := V[t][y]
        if p > prob || prob == 0 {prob, state = p, y}
    }
    return prob, path[state], idxs
}


func (v *HMM)Cut(src string) (ret []string) {
    _, path, idxs := v.Viterbi(src)
    // if only one rune, return it
    if len(idxs) == 1 {return []string{src}}
    last := 0
    for i, p := range path {
        switch {
        case p == "S":
            ret = append(ret, src[idxs[i]:idxs[i + 1]])
        case p == "B":
            last = idxs[i]
        case p == "E":
            ret = append(ret, src[last:idxs[i + 1]])
        }
    }
    return
}
