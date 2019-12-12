cyberstation-cli
================

JR の指定席の空席案内 ([JR サイバーステーション](http://www.jr.cyberstation.ne.jp/)) をコマンドラインで問い合わせます。

Go の標準ライブラリと `golang.org/x/text` だけでがんばります。

## Install

```
go get -u github.com/mikan/cyberstation-cli
```

## Usage

```
cyberstation -date YYYY/MM/DD -time HH:MM -from <乗車駅> -to <降車駅>
```

例:

```
cyberstation -date 2019/12/30 -time 22:00 -from 大垣 -to 品川
```

## Limitation

- 利用可能時間は 6:30 から 22:30 の間のみです
- 今のところ在来線のみに対応しています
- 駅名の候補が複数ある場合は未対応、完全一致する駅名をサイトで調べて渡してください

## License

サービスの利用に関する注意書きが以下の URL にあります:

> http://www.jr.cyberstation.ne.jp/license.html

本プログラムのソースコードは [BSD License](LICENSE) です。
