## 章立て
マニフェストの汎用化を行うOSS

## memo
### helm
* K8sのパッケージ（Chart）マネージャ（yumやaptに相当）、またはサービスカタログ
* Helmクライアントから要求を処理するサーバとして `Tiller Pod` の立ち上げが必要で `helm init` とかでいけちゃう
  * Contextを指定可能
  * kubectlと同じ認証情報を利用（`~/.kube/config`)
* クラスタに設定がない場合は下記のようにしてみた
```
$ kctl -n kube-system create serviceaccount tiller
serviceaccount/tiller created

$ kctl create --save-config clusterrolebinding tiller \
  --clusterrole=cluster-admin \
  --user="system:serviceaccount:kube-system:tiller"
clusterrolebinding.rbac.authorization.k8s.io/tiller created

$ helm init --service-account tiller
Error: unknown command "init" for "helm"
```
→ しかし、helmの3系からはこのtillerはなくてもよくなったみたいw
* Chartは [ここ](https://github.com/helm/charts)で確認可能
* リストの追加(incubatorも追加したければstableの代わりに指定)
```
$ helm repo add stable https://kubernetes-charts.storage.googleapis.com/

$ helm repo update
```
* 検索してインストール
```
$ helm search repo data
NAME                         	CHART VERSION	APP VERSION            	DESCRIPTION                                       
stable/datadog               	2.2.6        	7                      	Datadog Agent
・・・・・

$ helm inspect all stable/datadog
→ これで設定可能なパラメータなどがわかる

```
* 設定はクライアントだけで完結させるパターンとvaluesファイルの作成パターンとがある
```
■クライアント完結
$ helm install stable/datadog --version 2.2.6 --name datadog-chart-sample \
  --set datadog.apiKey=XXXXXXXXXXXX

■values.yaml
datadog:
  apiKey: XXXXXXXXXXX

$ helm install stable/datadog --version 2.2.6 --name datadog-chart-sample \
  --values values.yaml
```
* テンプレートからマニフェストを作成できる機能もあって便利
* Chartの中身はk8sのマニフェストテンプレートと、変数が主体
* 下記のようにcreateコマンドで雛形を作ることが推奨されている
```
$ helm create sample-charts

$ tree ./sample-charts
|-- Chart.yaml # Chartのメタデータ
|-- charts
|-- templates # マニフェストのテンプレート
|   |-- NOTES.txt # helm install時に出力するメッセージ
|   |-- _helpers.tpl # Releaseの名前を元に変数を定義
|   |-- deployment.yaml
|   |-- ingress.yaml
|   |-- service.yaml
`-- values.yaml # デフォルトの定義設定
```
* `requirements.yaml` で依存するChartを指定することも可能

### Ksonnet
* 複数のK8s環境に対して同じアプリケーションを環境の特殊設定を入れつつデプロイ可能
* 詳細は割愛

