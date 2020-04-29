## 章立て
マイクロサービスアーキテクチャとサービスメッシュ

## memo
主にサービスメッシュのお話。
基本的な仕組みは、Pod間の通信経路にプロキシを挟むことでトラフィックのモニタリングやトラフィックコントロールを行う。
モニタリングは、詳細なリクエストのerror rate、latency、connection count、request countのようなメトリクスを取得可能。
トラフィックのコントロールで言うと、DeploymentからDeploymentへのトラフィックを転送する際に転送先を柔軟にコントロールできる、
いわばカナリアリリースや、Cookieを使ったA/Bテストのようなことも可能。あとはCircuit Breakerなども可能。
他にも、Fault Injection、Rate Limit、Retry処理のようなことも実現可能です。

### Linkerd
* Scalaで書かれた15dと言うプロキシを使ってサービスメッシュを構成
* namedというものを用いると動的に転送ルールを設定可能

### Conduit (Linkerd v2)
* Rustで書かれたconduit-proxyを利用してサービスメッシュを構成
* Linkerdはマルチプラットフォーム向けに開発されたのに対して、ConduitはK8sに特化しているため、軽量化つシンプルに利用可能
* 現状はトラフィックコントロール機能はなく、モニタリング用途としてのサービスメッシュ機能が実装されている
* 各PodにSidecarとしてconduit-proxyが埋め込まれている形式
* [Getting Started](https://linkerd.io/2/getting-started/)で触ってみると面白い

### Istio
* C++で書かれたEnvoyを使ってサービスメッシュを構成、メトリクスはGrafanaを使って見ることが可能
* Service Graphによるマイクロサービス間のグラフ構造の可視化が可能（Zipkinによるグラフ化も可能）
* こちらも各PodにSidecarとしてEnvoy proxyが埋め込まれている形式
* トラフィックコントロールはEnvoy Proxyに転送ルールの設定が必要で、PilotコンポーネントがK8sのリソースで定義された情報をもとにEnvoyの設定を行う
* Envoyが集めたメトリクスはMixerコンポーネントが収集する
* マイクロサービス間の通信を `mutal TLS認証(mTLS)` で相互にTLS認証して暗号化することが可能
* [Quick Start](https://istio.io/docs/setup/getting-started/)を見てみると良い
* PodにEnvoy Proxyを内包するには下記のいずれか
  * `istiocl kube-inject`コマンドを利用
    * 既存のマニフェストをもとにEnvoy Proxyをインジェクトしたマニフェストを自動生成する
    * 標準出力を `kubectl apply -f` で渡してあげればファイルを新たに生成せずに適応可能
  * `mutating webhook admission controller`を利用して自動的にインジェクトする
  + namespace単位で自動インジェクションの有効化／無効化の設定が可能


