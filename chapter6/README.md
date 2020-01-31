## 章立て
Discovery&LBリソース

## memo
### Service
* 受信したトラフィックを複数のPodにロードバランシングする機能を提供
* k8sではサービスディスカバリはServiceが提供しているが主に次の3種類がある
  * 環境変数を利用したサービスディスカバリ
  * DNS Aレコードを利用したサービスディスカバリ
    * Serviceを作成した際に `[Service名].[NameSpace].svc.cluster.local` のようなFQDNが作られるとかもある
    * IPアドレスは極力Service名を使って名前解決すんべや
  * DNS SRVレコードを利用したサービスディスカバリ
    * SRVレコードは、Port名とProtocolを利用することで、サービスを提供しているPort番号を含めたエンドポイントをDNSで解決する仕組み
    * `[_ServiceのPort名].[_Protocol].[Service名].[NameSpace名].svc.cluster.local`
* L4ロードバランシングを提供

#### ClusterIP Service
* `type: ClusterIP` となる最も基本となるServiceで、k8sクラスタ内からのみ疎通できるInternal Networkの仮想IPが割り当てられる

#### ExternalIP Service
* **特定のK8s NodeのIPアドレス:Port** で受信したトラフィックをコンテナに転送する形で外部疎通性を確立するService
* これを利用するとk8sクラスタ外からも疎通が可能になる（Nodeが指定のPortでListen状態になるため）

#### NodePort Service
* **全てのK8s NodeのIPアドレス:Port** で受信したトラフィックをコンテナに転送する形で外部疎通性を確立するService
* こちらもk8sクラスタ外からも疎通が可能になる
* 利用できるポート範囲はk8sの多くの環境で30000〜32676（k8sのデフォルト値）※k8s masterで変更も可能
* Nodeを跨いだロードバランシングが行われるが、同じNode上のPodに転送させることも可能（Nodeに1つはあるDaemonSetのような例がそれに該当）
  * `spec.externalTrafficPolicy` に `Local` と設定するのみ（Clusterがdefault値）

#### LoadBalancer Service
* プロダクション環境でクラスタ外からトラフィックを受ける場合に、一番実用的で使い勝手が良いService
* この方は、クラスタ外のLoadBalancerに外部疎通性のある仮想IPを払い出すことが可能（AWSとかだとALBとかに払い出してあげる感じか）
* この方も同じNode上のPodに転送させることも可能

#### Headless Service（None）
* 対象となる個々のPodのIPアドレスが直接返ってくるService
* ロードバランシングするためのIPアドレスは提供されず、DNS Round Robinを使ってエンドポイントを提供

#### ExternalName Service
* 通常のServiceリソースとは異なり、Service名の名前解決に対して外部のドメイン宛のCNAMEを返す。別の名前を設定したい場合や、クラスタ内からのエンドポイントを切り替えやすくしたい場合など。
* 外部にあるサービスを利用する際に、この機能を使うと疎結合に保てるね！
  * 外部のエンドポイントをアプリケーション(Pod)には書かずに、このServiceを指定しておけば、あとはServiceの設定を変えればOK的な
* ClusterIP ServiceからExternalName Serviceに切り替える際は、 `spec.clusterIP` を明示的に空にする必要があることに注意

#### None-Selector Service
* k8s内に自由な宛先でここまで説明したロードバランシングを行える子（ClusterIPが大半みたい）
* ExternalName ServiceはCNAMEを返したのに対して、こいつはClusterIPのAレコードが返る（k8sのServiceの使用感のまま、外部むけにロードバランシングできる）

### ingress
* L7ロードバランシングを提供するリソース
* Network PolicyリソースにIngress/Egressという設定項目があるが、そのIngressとは無関係であることに注意
* よく使われるのは下記らしい。どのようなL7ロードバランサを作成するかは、このIngress Controllerによって異なる
  * GKE Ingress Controller
  * Nginx Ingress Controller
    * デプロイしたIngress Controllerはクラスタに登録された全てのIngressリソースを見てしまうことに注意
    * その時はIngress Classのアノテーションを付与すること
* AWSでは `AWS ALB Ingress Controller` なるものがあり、このリソースに該当しそう
* Ingressリソースの作成には事前準備が必要
  * Serviceをバックエンドとした転送を行う仕組みになっている
  * `type: NodePort` を指定すること
