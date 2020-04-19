## 章立て
ヘルスチェックとコンテナのライフサイクルについて

## memo
### ヘルスチェック
* k8sの標準のヘルスチェックは、Pod内のコンテナが実行しているプロセスが実行されているか
* その他にも下記の2つのヘルスチェック機構が用意されている
  * Liveness Probe: Podが正常に動作しているかの確認、失敗すればPodは再起動(podのrestartPolicyに沿う)
    * 例としては、長時間実行によりメモリリークで処理性能が低下など
  * Readiness Probe: Podがサービスインする準備ができているかの確認、失敗すればトラフィックを流さない（Podの再起動はしない）
    * 例としては、DBとの接続が完了しているか、キャッシュをロードし終わっているか
* Liveness ProbeもReadiness Probeも下記の3つのチェック方法が設定可能。どれか1つのコンテナでも失敗したらPod全体が失敗となる
  * exec: コマンドを実行し、終了コードが0でなければ失敗
  * httpGet: HTTP GETリクエストを実行し、status codeが200〜399でなければ失敗(pathやheaderの指定なども可能)
  * tcpSocket: TCPセッションが確立できなければ失敗
* 下記の5種類のヘルスチェック間隔が設定可能
  * initialDelaySeconds: 初回ヘルスチェック開始までの遅延
  * periodSeconds: ヘルスチェックの間隔
  * timeoutSeconds: タイムアウトまでの秒数
  * successThreshold: 成功と判断するまでのチェック回数
  * failureThreshold: 失敗と判断するまでのチェック回数

### コンテナのライフサイクル
* コンテナのプロセス停止 or ヘルスチェックの失敗時にコンテナの再起動するかどうかは `spec.restartPolicy` で指定可能
  * Always: Podが停止すると、常にPodを再起動させる（デフォルトはこれ）
  * OnFailure: Podの停止が予期せぬ停止（終了コードが0以外）の場合、Podを再起動させる
  * Never: Podが停止しても、Podを再起動しない
* Init Containerとは、Pod内のメインとなるコンテナを起動する前に別のコンテナを起動させるための機能（複数指定可能で、上から1つずつ実行される）
  * `kubectl get pods --watch` のstatusに進捗が表示される
* コンテナの起動後（postStart）と、停止直前（preStop）に任意のコマンドを実行可能
  * postStartは非同期に行われるので、entrypointが実行されるより前に実行される可能性があるので注意
* Podの削除が行われる際には、非同期に「preStop処理+SIGTERM」と「Serviceからの除外設定」が行われる
  * ただし厳密には同時ではないので、Serviceの例外処理が終わるであろう数秒間を「preStop処理+SIGTERM」で待機するか、そこでリクエスト受けつつ停止処理を行うか
  * terminationGracePeriodSecondsには「preStop処理+SIGTERM」が完了するのに十分な時間を確保すること（デフォルトは30秒）
    * preStop処理でterminationGracePeriodSecondsを使い切ってしまったら、2秒だけSIGTERM処理の時間が確保されるみたい
  * kubectlコマンドで削除するときは「--grace-period」オプションを利用して作成済みのPodに対して適応することも可能
