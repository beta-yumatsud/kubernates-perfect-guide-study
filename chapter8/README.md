## 章立て
Config & Storageリソースについて

## memo
### Clusterリソース
* セキュリティ周りの設定やクォータ設定など、クラスタの挙動を制御するためのリソース
* Node, NameSpace, PersistentVolume, ResourceQuota, ServiceAccount, Role, ClusterRole, RoleBinding, ClusterRoleBinding, NetworkPolicyなどがある
* Node
  * `status.addresses` にIPアドレスとかホスト名が保持されている
  * `status.allocatable` と `status.capacity` にnodeのリソース状況の記載がある
  * nodeのstatusがreadyではない場合、`status.conditions`をみると原因がわかる
  * `status`を参照すると良いね
* NameSpace
  * 仮想的なクラスタの分離機能
  * マニフェストファイルでもkubectlコマンドでも作成可能

### Metadataリソース
* クラスタ上にコンテナを起動させるのに利用するリソース
* LimitRange, HorizontalPodAutoScaler, PodDisruptionBudget, CustomResourceDefinitionなどがある
