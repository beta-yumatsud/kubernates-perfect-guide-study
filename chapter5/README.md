## 章立て
Workerloadsリソースについて

## memo
### Pod
* Podのデザインパターンは大きく分けて下記の3種類がある
  * サイドカーパターン: メインコンテナに機能を追加する
  * アンバサダーパターン: 外部システムとのやり取りの代理を行う
  * 外部からのアクセスのインターフェースとなる
* `kubectl exec -it <<pod>> <<command>>` はローカルで試す時にはよく使いそう
* dockerのENTRYPOINT/CMDとcommand/args
  * ENTRYPOINTとcommandが対応関係
  * CMDとargsが対応関係
* Pod名の制限（RFC1123）
  * 利用可能文字: 英数字の小文字と`-`と`.`（`_`は使えないぞよ）
  * 始まりと終わりの文字は英小文字
* Pod内の全コンテナの`/etc/hosts`を書き換える機能があるんだってさ（spec.host.Aliases）

### ReplicaSet
* ReplicaControllerというのが最初は使われていたようだが、歴史的背景もあり、今後は使われなくなるとのこと（アルヨネー）
* ReplicaSetはk8sがPodの監視を行うことでPodの数を調整しますが、その際、特定のラベルが付けられたPodの数をカウントする形で実現しているらしいよ
  * `selector.matchLabels`の部分で指定しているものがまさにそれで、`template.metadata.labels`の部分に対応する
  * その2つが一致しないとエラーで作成ができないか、レプリカ数を増やそうとPodが永遠に作られ続けてしまうかも（こわっ）
  * ReplicaSetをスケールさせる方法は下記の2つ。基本的には前者を使いませう（コード管理のためにも）
    * マニフェストファイルを書き換えて `kubectl apply -f` を実行
    * `kubectl scale rs <rs名> --replicas <スケールさせる数>`コマンドを利用

### Deployment
* 1つのPodでもDeploymentを使うようにしておくと、セルフヒーリングもローリングアップデートも使えるから良い
* apply時に`--record`を付けておくとアップデートを行った際の履歴を保持できる
  * `kubectl set image deployment <deployment名> <コンテナ名=コンテナ:新バージョン>`とかで更新できる
  * `kubectl rollout status`コマンドでアップデート状況を確認できる
* 新しいReplicaSetが作られる条件は `作成されるPodの内容の変更` が必要
  * `spec.template以下の構造体のハッシュ値` を計算して、それを利用したラベルづけで管理しているから
  * ReplicaSetの `spec.template.metadata.labels.pod-template-hash`がまさにそれ
* 変更履歴は `kubectl rollout history deploy <deployment名>` で確認できる
  * `--record`を付けていないとCHANTE-CAUSE部分は空になる
  * `--revision`オプションで該当リビジジョンの詳細を取得できるよ
* ロールバックは `kubectl rollout undo` を利用する
  * `--to-revision` でどこまでロールバックしたいかを指定可能
  * 何も指定しないと `--to-revision 0` と同じで、前のリビジョンに戻る
* 即時適用を一時停止するには `kubectl rollout paust` , 再開は `kubectl rollout resume`
* RollingUpdateでは `maxUnavailable(アップデート中に許容される不足Pod数)` と `maxSurge(超過Pod数)`を指定することでリソース消費量やアップデートのスピードなどを制御する
  * `spec.strategy.rollingUpdate`いかに設定を行う
  * パーセンテージの指定も可能。デフォルトではそれぞれ25%指定になっている
* アップデート時には下記のような他のパラメータを指定可能
  * `minReadySeconds`: 次のPodの入れ替えが可能と判断するまでの最低秒数
  * `revisionHistoryLimit`: ロールバック可能な履歴数（保存するReplicaSetの数）
  * `progressDeadlineSeconds`: タイムアウト時間。これが経過した場合、自動でロールバックされる

### DaemonSet
* 各NodeにPodを1つずつ配置するリソース！（ゆえにレプリカ数などの指定はできない）
  * とはいこのNodeには配置しない！などの指定は可能
  * 用途としては、ログをホスト単位で収集するFluentdや、各Podのリソース使用状況やNodeの状態をモニタリングするDatadogなど、全Nodeで必ず動作させたいプロセスのために利用することが多い
* アップデート戦略には `RollingUpdate` と `OnDelete（次回再作成時や手動更新）`の2つがある

### Job
* コンテナを利用して一度限りの処理を実行させるリソース（N並列で実行しながら、指定した回数のコンテナを実行:正常終了を保証するリソース）
* labelやselectorを明示的に付与することも可能だが、k8sがユニークなuuidを自動生成するため明示的に付与することは推奨されていない
* restartPolicyには下記の2つがある
  * `OnFailure`: 再度同一のPodを利用してJobを再開
  * `Never`: Pod障害時には新規のPodが作成される
* completions, parallelism, backoffLimitの組み合わせ例
  * 1回だけ実行するタスクは`completions:1`, `parallelism:1`, `backoffLimit:1`
  * N並列で実行させるタスク(並列ジョブ的な)は`completions:M`, `parallelism: N`, `backoffLimit: P`
  * 1個ずつ実行するワークキュー(コンテナを起動しては特定のタスクを実行し続ける状態)は`completions:未指定`, `parallelism: 1`, `backoffLimit: P（ただしデフォルトは6で無限とかは無理ぽ）`
  * N並列で実行するワークキューは`completions:未指定`, `parallelism: N`, `backoffLimit: P（ただしデフォルトは6で無限とかは無理ぽ）`

### CronJob
* Cronのようにスケジュールされた時間にJobを作成するのがCronJob
* CronJobがJobを管理し、JobがPodを管理するという3層の親子構造
* 一時停止は `spec.suspend: true` にする
  * `kubectl patch cronjob sample-cronjob -p '{"suspend":true}'` とかでpatchを当てることも可能
* 前のjobが正常完了していれば、指定時間に再度実行されるが、失敗していた場合の実行ポリシーは `spec.concurrencyPolicy` で下記のように指定可能
  * `Allow`: 同時実行に対して制限を行わない
  * `Forbid`: 前のJobが終了していない場合、次のJobは実行しない（同時実行は行わない）
  * `Replace`: 前のJobをキャンセルし、Jobを開始する
* `spec.startingDeadlineSeconds`は開始時刻がどこまで遅れて良いかを制限するもの
* `spec.successfulJobsHistoryLimit`は成功したJobの数を保存しておき、`spec.failedJobsHistoryLimit`は失敗したJobの数を保存しておく（いずれもデフォルト値は3）
