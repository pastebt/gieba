package gieba

import "testing"
import "strings"


func TestIndex(tst *testing.T) {
    src := "测试一下"
    idxs := make([]int, 0, len(src))
    for i := range src {
        idxs = append(idxs, i)
    }
    idxs = append(idxs, len(src))
    for I, i := range idxs {
        for _, j := range idxs[I+1:] {
            tst.Log(src[i:j])
        }
    }
}


func TestCutNoHMM(tst *testing.T) {
    f := FreqNew()
    f.LoadFile("zh_dict.txt")
    f.hmm = HMMNew("/usr/share/gieba/data/")
    okdata := [][]string {
        {"我来到北京清华大学", "我/ 来到/ 北京/ 清华大学"},
        {"他来到了网易杭研大厦", "他/ 来到/ 了/ 网易/ 杭研/ 大厦"},
    }
    for _, dat := range okdata {
        ret := f.Cut(dat[0])
        res := strings.Join(ret, "/ ")
        if res != dat[1] {
           tst.Error(ret, dat)
        }
    }
}
