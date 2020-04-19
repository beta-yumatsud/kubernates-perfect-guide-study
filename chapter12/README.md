## 章立て
高度で柔軟なスケジューリング

## memo
* スケジューリングでは下記の考え方をする
  * Affinity: 特定の条件に一致するところにスケジューリングする
  * Anti-Affinity: 特定の条件に一致しないところにスケジューリングする
* k8sのスケジューリングには大きく分けて下記の2つがある
  * Podのスケジューリング時に特定のNodeを選択する方法、さらに下記の5種類がある
    * nodeSelector(Simplest Node Affinity): 簡易的なNode Affinity機能
    * Node Affinity: 特定のノード上だけで実行
    * Node Anti-Affinity: 特定のノード以外で実行
    * Inter-Pod Affinity: 特定のPodがいるドメイン（ノード、ゾーン、etc）上で実行する
    * Inter-Pod Anti-Affinity: 特定のPodがいないドメイン（ノード、ゾーン、etc）上で実行する
  * Nodeに対して汚れ（Taints）をつけておき、それを許容（Tolerations）できるPodのみスケジューリングを許可する方法
* ビルトインノードラベルとはk8sで付与するラベル（ホスト名/OS/アーキテクチャなど）があり、これに追加で付与することも可能（kubectl labelコマンド）
* nodeSelectorは指定のラベルを持つノートに配置
  * ex. 高いディスクIO性能が必要なPodを配置したい場合など
  * spec.nodeSelectorにラベルを指定すればOK
* Node Affinityは spec.affinity.nodeAffinity に記載。必須のスケジューリングポリシーと優先スケジューリングポリシーの設定がある
  * 詳しくは [sample](sample-node-affinity.yaml)を参照
  * 条件に合致するNodeが存在しない場合は、Podはスケジューリングされない
* matchExpressionsのオペレータとset-based条件
  * 利用可能な operator は下記の6つ
    * In: `A IN [B...]`、ラベルA(ex. key）の値が[B...]（ex. values）の中のいずれか1つ以上と一致する
    * NotIn: `A NotIn [B...]`、ラベルAの値が[B...]の中のいずれにも一致しない
    * Exists: `A Exists []`、ラベルAが存在する（ex. keyに指定したものがあればOK的な）
    * DoesNotExists: `A DoesNotExists []`、ラベルAが存在しない
    * Gt: `A Gt [B]`、ラベルAの値がBよりも大きい
    * Lt: `A Lt [B]`、ラベルAの値がBよりも小さい
  * これらは ReplicaSet/Deployment/DaemonSet/StatefulSet/JobのSelectorで指定可能
* Node Anti-Affinity（実際は存在しない）は spec.affinity.nodeAffinity の条件で否定系のオペレータを指定するだけ
* Inter-Pod AffinityはPod定義内の spec.affinity.podAffinity に条件記述することで Pod Affinity を設定する
  * Pod同士を近づけて配置することで、Pod間通信のレイテンシを下げることが目的だったりする
  * failure-domain.beta.kubernetes.io/zoneに指定すると同じラベルを持つPodと同じソーンに配置
* Inter-Pod Anti-Affinityは spec.affinity.podAntiAffinity で指定可能、設定条件などはInter-Pod Affinityと同様
* TaintsとTolerations
  * Pod側が提示して、Node側が許可するイメージ
  * 条件に合致しないPodはNode上から追い出すことも可能（Node Affinityとかは残っちゃう）
  * 利用シーンとしては、対象のノードを特定の用途にむけた専用ノードにする場合など（productionとかGPUとか）
  * `kubectl taints`で `Key=Value:Effect` の書式で実行。Effectはマッチしない時の挙動で、下記の3つがあり、この順で条件が厳しくなる
    * PreferNoSchedule: 可能な限りスケジューリングしない
    * NoSchedule: スケジューリングしない（すでにスケジューリングされているPodはそのまま）
    * NoExecute: 実行を許可しない（すでにスケジューリングされているPodは停止される）
  * Tolerations(許容)はKeyの条件式とEffectを指定し、Taints(汚れ)で付与されたKeyとEffectの両方が同じだった場合に条件が一致したと判断して許容する
    * こちらは `spec.tolerations` で指定
* PriorityClassによるPodの優先度と退避が可能
 