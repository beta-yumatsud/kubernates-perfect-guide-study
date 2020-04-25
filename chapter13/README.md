## 章立
セキュリティ（広いな）

## memo
### Service Account
* K8sにはユーザに似たコンセプトとして `UserAccoun` と `SeviceAccount` がある
  * UserAccount
    * EKSとかでいうとIAMとリンクするようになってるみたい
    * K8sの管理対象ではない
    * クラスタレベルでの存在でNameSpaceの影響を受けない
  * ServiceAccount
    * K8sの世界の中だけで完結するもの
    * Podで実行されるプロセスのために割り当てられるもので、ServiceAccountベースで認証/認可を実施
      * 指定しないと default ServiceAccount が割り当てられる
    * NameSpaceに紐づく
* `kubectl create serviceaccount` やyamlを作成して作成可能
  * ServiceAccountを作成すると、Secretsが作成されるが、Secretsは基本的に自動で作成されるもの
  * Secretはトークンと証明書から構成される
  * このトークンを元にしてK8s APIへの認証情報として利用することが可能
* Podに `spec.servceAccountName` で明示的にServiceAccountを指定可能
  * そうするとトークンがVolumeとして自動的にマウントされる
  * `automountSeviceAccountMount` を `false` で指定すると自動マウントはしなくても良い
* `kunernetes/client-go` を使ってK8s APIを使えたりもする

### RBAC
* Role Based Access Controlとはどういう操作を許可するのかを定めたRoleを作成し、Service AccountなどのUserに対してRoleを紐付ける（RoleBinding）ことで管理する
  * K8s 1.9から利用可能になった `AggregationRule` を利用することでRoleが別のRoleを参照可能
  * 集約した側のruleは無視されることに注意
  * これをやっておくと、集約された側の変更も、集約する側に自動的に反映されるみたい
* RBACにはNamespaceレベルのリソースとClusterレベルのリソースの2種類が用意されている
* RoleとClusterRole
  * いくつかのプリセットが用意されているので、ざっくりとしたものはこれを使うと良い
  * Role作成時に対象のリソースに対して実行できる操作は`rules[].verbs`で指定（create、delete、getなど）
* sampleのyamlをみてみた方がイメージ湧くかも、あとEKSとかを使うだけならあんまり意識しない可能性もある

### Security Context
* 個々のコンテナに対するセキュリティ設定（設定可能項目はググろうぞ）
* `spec.containers[].securityContext.hogehoge` で設定が可能
* rootファイルシステムをRead Onlyにするなどの設定はしておいても良いのかなと思う

### Pod Security Context
* Pod（全てのコンテナ）に対するセキュリティ設定（設定可能項目はググろうぞ）
* `spec.securityContext.hogehoge` で設定
* sysctlによるカーネルパラメータの設定とかは面白い（が、クラウド利用でどこまでやるのだろうか）

### Pod Security Policy
* K8sクラスタに対してセキュリティポリシーによる制限を行う機能
* ex. SecurityContextなどで設定した項目のデフォルト設定値や制限を行うことが可能（設定可能項目はぐぐろうぞ）

### Network Policy
* K8sクラスタ内でPod同士が通信する際のトラフィックルールを規程できるもの
  * 規程しないと全てのPod同士がクラスタ内で通信可能
  * 設定すれば、NameSpace毎にトラフィックを転送しないようにしたり、特定のPod間だけ通信可能とかにもできる
  * ただし全てのK8s環境で使えるわけではないので要注意
  * EKSは [この手順](https://docs.aws.amazon.com/ja_jp/eks/latest/userguide/calico.html) で設定可能
* この子はIngressとEgressから成り立つ

### 認証/認可とadmin control
* admin controlはK8s APIサーバにリクエスト制御を追加で行う仕組み
* K8sでは下記のフェーズを通ってリソースが登録される
  * Authentication/AuthN(認証): 正しいユーザであるか
  * Authorization/AuthZ(認可): 権限があるか
  * Admin Control: リソースの改変や柔軟なチェックが可能
* 基本的にはクラスタ管理者が設定したりするものみたいだよ
* PodPreset: Podを作成する際にデフォルト設定を埋め込むリソース（特定のラベルつきPodに自動的に環境変数を埋め込むなど）
  * env, envFrom, volumes, volumeMountsを埋め込み可能
  * デフォルトでは有効になっていないため、有効化が必要
  * ただし、PodPresetの設定とPodの設定が1つでもかぶっていると全てのPodPresetが設定されないので注意
  * Pod側でPodPresetの対象から除外されたい場合の設定とかは可能だよ