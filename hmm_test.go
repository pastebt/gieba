package gieba


import "testing"

import (
    "os"
    "log"
    "strings"
    "io/ioutil"
    "encoding/json"
)


func TestLoadJson(tst *testing.T) {
    infn := "/usr/share/gieba/data/prob_emit.json"
    //infn := "/home/fure/golang/fgdwcfgo/segm/data/CN/prob_start.json"
    //infn := "/home/fure/golang/fgdwcfgo/segm/data/CN/prob_trans.json"
    fin, err := os.Open(infn)
    if err != nil {
        log.Fatal(err)
    }
    defer fin.Close()
    dat, err := ioutil.ReadAll(fin)
    emit := make(map[string]map[string]float64)
    err = json.Unmarshal(dat, &emit)
    //start := make(map[string]float64)
    //err = json.Unmarshal(dat, &start)
    //trans := make(map[string]map[string]float64)
    //err = json.Unmarshal(dat, &trans)
    if err != nil {
        tst.Error(err)
    }
}


func TestHMMNew(tst *testing.T) {
    v := HMMNew("")
    if v == nil {
        tst.Error("Failed to new vitebi")
    }
    v.LoadJsons("/usr/share/gieba/data/")
    _, p, i := v.Viterbi("测试一下")
    show("path=", p)
    show("idx=", i)
}


func TestHMMCut(tst *testing.T) {
    v := HMMNew("/usr/share/gieba/data/")
    okdata := [][]string {
        {"测试", "测试"},
        {"测试一下", "测试/ 一下"},
        //{"他来到了网易杭研大厦", "他/ 来到/ 了/ 网易/ 杭研/ 大厦"},
        //{"我来到北京清华大学", ""},
        //{"小明硕士毕业于中国科学院计算所，后在日本京都大学深造", ""},
        }
    for _, dat := range okdata {
        r := v.Cut(dat[0])
        show("ret=", r)
        s := strings.Join(r, "/ ")
        if s != dat[1] {
            tst.Error(r, dat)
        }
    }
}
