gieba
=====

jieba in golang

jieba url: https://github.com/fxsjy/jieba

cd $GOROOT/src/pkg
mkdir -p github.com/pastebt
cd github.com/pastebt
git clone https://github.com/pastebt/gieba
go test github.com/pastebt/gieba

mkdir /usr/share/gieba
cd /usr/share/gieba
ln -s $GOROOT/src/pkg/github.com/pastebt/gieba/data/ .
