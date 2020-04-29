## 章立て
K8s環境でのCI/CD

## memo
### Spinnaker
* Netflix作成のOSS
* 様々なデプロイ方法（カナリア、ブルグリなど）、様々なCIツールからのイベントトリガー、Slackやメールでの通知など
* 環境構築には `Halyard` を用いて実施することが推奨されている

### Skaffold
* Googleが開発したOSSで、K8s向けにビルドとデプロイを自動化するもの
* ビルドしたイメージをレジストリにPushせずにK8sクラスタへデプロイすることが可能
* マニフェストファイルの形式は通常の形式と同じでOKで、ファイル名のプレフィックスに「skaffold-」をつけること
* 下記コマンドで開発モードで起動
```
$ skaffold dev
```
* hot reloadのような機能も持っていて、コード修正するとimageのbuild、push、deployまで行ってくれる
  * ん〜、開発中はpushとかまではいらんなw
  * そんな時は、ローカルに限り、imageのpushを省くことも可能
  * `build.local.skipPush=true` をマニフェストに指定すればOK
* skaffoldの設定ファイルに `profiles` を指定することで、各環境ごとの設定をすることが可能
  * `skaffold dev --profile devProfile`
* `skaffold dev`は基本的に開発用途で、本番で使う場合は `skaffold run` を使うこと
  * これは一回だけ実行されるので、パイプラインを1回だけ実行するために使えたりする
