cyberstation-cli
================

JR の指定席の空席案内 ([JR サイバーステーション](http://www.jr.cyberstation.ne.jp/)) をコマンドラインで問い合わせます。

Go の標準ライブラリと `golang.org/x/text` だけでがんばります。

## Install

```
go get -u github.com/mikan/cyberstation-cli/cmd/cyberstation
```

## Usage

```
cyberstation -date YYYY/MM/DD -time HH:MM -from <乗車駅> -to <降車駅> -group <列車種別>
```

在来線の例:

```
cyberstation -date 2019/12/30 -time 10:00 -from 品川 -to 大垣
```

新幹線の例:

```
cyberstation -date 2019/12/30 -time 10:00 -from かみのやま温泉 -to さくらんぼ東根 -group 3
```

`-group` の数値は以下のルールに従って指定する必要があります:

- 1: 東海道・山陽・九州新幹線 (のぞみ, ひかり, みずほ, さくら, つばめ)
- 2: 東海道・山陽新幹線 (こだま)
- 3: 東北・北海道・山形・秋田新幹線 (はやぶさ, はやて, やまびこ, なすの, つばさ, こまち)
- 4: 上越・北陸新幹線 (とき, たにがわ, かがやき, はくたか, あさま, つるぎ)

## Limitation

- 利用可能時間は 6:30 から 22:30 の間のみです
- 駅名の候補が複数ある場合は未対応、完全一致する駅名をサイトで調べて渡してください

## License

サービスの利用に関する注意書きが以下の URL にあります:

> http://www.jr.cyberstation.ne.jp/license.html

本プログラムのソースコードは [BSD License](LICENSE) です。
