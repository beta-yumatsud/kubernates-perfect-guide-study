## 章立て
Config & Storageリソースについて

## memo

### Environment
* K8sで環境変数を渡す際は、PodテンプレートにenvまたはenvFormを指定するが、大きく分けて下記の5つの情報源から環境変数を埋め込むことが可能
  * 静的設定
    * `spec.containers[].evn` に設定
    * タイムゾーンの設定とかもいけちゃうぜよ
  * Podの情報
    * どのノードで起動しているかや起動時間などのPodに関する情報は `fieldRef` で参照可能
  * コンテナの情報
    * resourceFieldRefでコンテナの情報は参照可能
  * Secretリソースの情報
    * 機密情報などは別途Secretリソースを作成し、環境変数として参照させることが推奨されている
  * ConfigMapリソースの設定値
    * key-value 形式の設定ファイルはこれを使うと良い
* 環境変数利用上の注意
  * K8sではcommandやargsで実行する際に環境変数を利用できない＞＜
  * 厳密にはPodマニフェストないでは、`$()` で参照可能（`${}`だとエスケープされちゃうんだっさ。かつ、OS設定での環境変数などはむりぽ、マニフェストないで定義したもののみ参照可能）

### Secret
* DBのユーザ名やパスワードなどの機密情報を扱う時に利用するリソース
  * ただし、Secretを定義しているマニフェストは、Base64エンコードされているだけで、暗号化はしていない
  * このような問題を解決すべく、マニフェストを暗号化する `kubesec` などがある（AWS KMSなどを利用することで簡単に `data.*` 部分を暗号化できるらしい）
* 分類
  * Generic（type: Opaque）→DBのユーザ名やパスワードはこれを使うおー
  * TLS（type: kubernetes.io/dockercongigjson）
  * Dockerレジストリ（type: kubernetes.io/dockercongigjson）
  * Service Account（type: kubernetes.io/service-account-token）
* Generic
  * 作成方法は下記の2つ
    * `--from-file`
      * ファイル名がそのままkeyとなるので、拡張子などは外しておくの良いみたい
      * あとはファイルに改行コードが入れないように注意
      * データ自体は、base64エンコードされている
    * `--from-env-file`
      * 1つのファイルから一括で作成するのに便利
      * 下記のようなファイルを用意した場合は `kubectl create secret generic --save-config sample-db-auth --from-env-file ./env-secret.txt` で作成できる
      ```env-secret.txt
      username=hogehoge
      password=hogehogepassword
      ```
    * `--from-literal`
      * `--from-literal=key=value` で指定すればOK
    * `-f` ( マニフェストから作成 )
      * base64エンコードした値をマニフェストに埋め込むことに注意
  * 1つのSecretに格納可能なデータサイズは1MB
* TLS
  * 証明書として利用するsecretはtlsタイプを指定
  * Ingressとかから利用することが一般的で、秘密鍵と証明書のファイルからSecretを作成することが望ましい（`--key`, `--cert`）
* Dockerレジストリ
  * プライペーとリポジトリを利用している場合はDockerイメージの取得の際には認証が必要でそれ用
  * typeにはdocker-registryを指定で、kubectlから直接作成するのが良いみたい
  ```
  $ kubectl create secret --save-config docker-registry sample-registry-auth \
    --docker-server=REGISTRY_SERVER \
    --docker-username=REGISTRY_USER \
    --docker-password=REGISTRY_USER_PASSWORD \
    --docker-email=REGISTRY_USER_EMAIL
  ```
  * イメージの取得のためPod定義の `spec.imagePullSecrets` にdocker-registryタイプのsecretを指定
* secretをコンテナから利用する場合は、大きく分けて下記の2パターン
  * 環境変数として渡す
    * Secretの特定のkeyのみ
      * `spec.containers[].env` の `valueFrom.secretKeyRef` を使って指定
    * Secretの全てのkey
      * `spec.containers[].envFrom` の `secretRef` を使って指定
  * Volumeとしてマウントする
    * Secretの特定のkeyのみ
      * `spec.volumes[]` の `secret.items[]` を使って指定
    * Secretの全てのkey
      * `spec.volumes[]` の `secret` を使って指定
* kubesec
  * K8sのsecretを安全に管理するためのOSS（AWSの場合KMSを利用してSecretを暗号化を行うことが可能
  * ファイル全体ではなく、Secretの構造（data.＊）を保ったまま値だけ暗号化する、故にGitとかにもあげれちゃう
  * macOSであれば `brew install shyiko/kubesec/kubesec` で入れれるおー

### ConfigMap
* 設定情報などのkey-valueで保持できるデータを保存しておくリソース
  * nginx.configやhttpd.configのような設定ファイルそのもの自体も保存可能！
* 作成方法はGenericタイプで下記の3通り
  * `--from-file`
    * こちらもファイル名がkey名になっていく（ `--from-file=nginx.conf=sample-nginx.conf` とかで変更可）
    * `kubectl create configmap --save-config sample-configmap --from-file=./nginx.conf` などなど
    * `kubectl get configmap <configmap名>` で確認できるし `describe` で詳細も見れる
  * `--from-literal`
    * 直接key-valueを指定して作成
    * `kubectl create configmap --save-config web-config --from-literal=connection.max=100 --from-literal=connection.min=10` 的な
  * `-f`（マニフェストファイル）
* 1つのConfigMapにつき格納可能なデータサイズは1MB
* ConfigMapをコンテナから利用する場合は、Secretと同じく大きく分けて下記の2パターン
  * 環境変数として渡す
    * Secretの特定のkeyのみ
      * `spec.containers[].env` の `valueFrom.configMapKeyRef` を使って指定
    * Secretの全てのkey
      * `spec.containers[].envFrom` の `configMapRef` を使って指定
  * Volumeとしてマウントする（複数行にわたるnginx.confとかで使う）
    * Secretの特定のkeyのみ
      * `spec.volumes[]` の `configMap.items[]` を使って指定
    * Secretの全てのkey
      * `spec.volumes[]` の `configMap` を使って指定

### 永続化領域
* そもそもVolume、PersistentVolume、PersistentVolumeClaimのの違い
  * Volumeはあらかじめ用意された利用可能なボリュームをマニフェストに指定して利用可能にするもの（故にVolume自体の作成や変更などの操作はできない
  * PersistentVolumeは外部の永続ボリュームを提供するシステムと連携して、ボリュームの作成や削除などの操作を行うもの（マニフェストなどからリソースを別途作成するイメージ）
  * PersistentVolumeClaimはPersistentVolumeリソースの中からアサインするためのリソース
    * PersistentVolumeはクラスタにボリュームを登録するだけなので、実際にPodから利用するにはこれを定義してい利用する必要がある
    * Dynamic Provisioning機能は、PersistentVolumeClaimを利用するタイミングで、PersistentVolumeを動的に作成可能


### Volume
* 下記のような様々なプラグインを用意
  * emptyDir
    * Pod用の一時的なディスク領域として利用可能
  * hostPath
    * Node上の任意の領域をコンテナにマッピングする
  * downwardAPI
    * Podの情報などをファイルとして配置する
  * projected
    * Secret/ConfigMap/downwardAPI/serviceAccountTokenのボリュームマウントを1箇所のディレクトリにに集約する
  * nfs
  * iscsi
  * cephfs

### PersistentVolume(PV)
* ネットワーク越しにディスクをアタッチするタイプのディスクになるよう
* 作成する際には下記のような項目の設定が可能
  * ラベル
    * 割り当てを行う際に目印となるものを付けておこう（type, envなどなど）
  * 容量
    * Dynamic Provisioningを利用できない環境では小さめのPersistentVolumeを用意しておくこと
  * アクセスモード
  * ReadWriteOnce(RWO)、ReadWriteMany(RWM、GCPやAWSなどでは許可されていない)、ReadOnlyMany(ROM)などがある
  * Reclaim Policy
    * 利用し終わったあとの処理方法を制御するポリシー、`spec.persistentVolumeReclaimPolicy`を下記の3つから指定可能
      * Delete
        * 実態の削除
        * GCPやAWSなどで確保される外部ボリュームのDynamic Provisioning時に利用されることが多い
      * Retain
        * 実際を消さずに保持
        * 他の既存のPersistentVolumeClaimによって再度マウントされることはないみたい（新規ではあり）
      * Recycle
        * PersistentVolumeのデータを削除し、再利用可能な状態にする
        * 他のPersistentVolumeClaimによって再度マウントされえる
        * ただし、将来のK8sにて廃止が検討されているので、Dynamic Provisioningを利用すること
  * マウントオプション
  * StorageClass
  * 各プラグインに特有の設定

### PersistentVolumeClaim(PVC)
* 永続化領域の要求を行うリソース
  * PersistentVolumeはPVC経由で利用する形になる
* PVCを作成する際には下記のような項目が設定可能で、いずれもPersistentVolumeで定義した値。
  * ラベルセレクタ
  * 容量
  * アクセスモード
  * StorageClass
    * ユーザがPVCを使って Dynamic Provisioning でPVを要求する際にどういったディスクが欲しいのかを指定するために利用される
* 注意としては、PVCの容量が、PersistentVolumeの容量よりも小さい小さければ割り当てられてしまう点（要は容量にあまりが出ちゃって無駄が発生しちゃうよーということ）
* PodからPVCを利用するには `spec.volumes` に `persistentVolumeClaim.claimName` を指定する
* Dynamic Provisioning
  * これを使ったPVCの場合、PVCが発行されたタイミングで動的にPersistentVolumeを作成して割り当てる
  * これにより「事前にPersistentVolumeを作成しておく必要がない」「容量の無駄が生じない」というメリットが出てくる
  * ただし、多くのProvisionerが一部のアクセスモードしかサポートしていない点は注意す
  * リサイズがサポートされいているボリュームプラグインを使用している時はPVCの拡張（リサイズ）を行うことが可能（PersistentVolumeClaimResize）
* volumeMountsで利用可能なオプション
  * readOnlyマウント（hostPathとかでホスト上の領域を見せる必要がある時は最低限これを指定するなどなど）
  * subPath（特定のディレクトリをルートとしてマウントする機能）
