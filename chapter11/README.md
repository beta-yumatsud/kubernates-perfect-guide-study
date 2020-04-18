## 章立て
メンテナンスとノードの停止について

## memo
* k8sのノードは `SchedulingEnabled` と `SchedulingDisabled` のどちらかのステータス
  * SchedulingDisabledになると、スケジューリングの対象から外される
  * SchedulingDisabledにするには `kubectl cordon` を利用
* ノード上の全てのpodを排出処理するには `kubectl drain` を使う
* PodDisruptionBudgetを設定すると条件にマッチするPodの最小の起動数と最大の起動数をみながらNodeからのPod追い出しを行うことが可能