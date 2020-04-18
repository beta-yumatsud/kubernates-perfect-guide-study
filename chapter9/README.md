## 章立て
リソース管理とオートスケーリングについて

## memo
* リソースの種別と単位は下記
  * CPU: `1 = 1000m = 1 vCPU`
  * memory: `1G = 1000M(1Gi = 1024Mi)`
  * `requests`はリソースの下限、ただしこれすらリソースに空きがないときはスケジューリングされない
  * `limits`は上限を示すが、あまりにも大きいものを指定するとクラスタ全体の負荷が上がるので要注意
  * Device Plugins機能がサポートされたので、GPUとかも可能
* cluster autoscalerなるもので、クラスタ自体のオートスケーリングが可能（Nodeを自動的に追加的なね）
  * requestsなどを適切に設定していないと、実際の負荷がそうでもないにスケールされちゃうし、逆に負荷が高いのにスケールされないなんてこともある
  * requrestとlimitsに顕著な差をつけない、requestsを大きくしすぎない
* LimitRangeによるリソース制限が可能
  * Namespaceに対して制限をかけるため、Namespaceごとに設定する必要あり（新規で作成されるPodにのみ適応される）
  * 制限可能な項目
    * default: デフォルトのlimits
    * defaultRequest: デフォルトのRequests
    * max: 最大リソース
    * min: 最小リソース
    * maxLimitRequestRatio: Limits/Requestsの割合
  * 設定可能なリソースは Pod(max/min/maxLimitRequestRatio), Container, PersistentVolumeClaim(max/min)
* Podに設定されるQoS Classはrequestsやlimitsから自動でk8sが設定してくれる
  * Podの `status.qosClass` で確認可能
  * BestEffort, Guaranteed, Burstableの3つがあるよ
* ResourceQuotaを使うと各NameSpaceごとの利用可能なリソースを制限できる(使用量にも制限がかけれる見たい)
* HorizontalPodAutoscaler(HPA)は、Deployment/ReplicaSetのレプリカ数をCPU負荷などに応じて自動的にスケールさせるリソース
  * PodにResource Requestが設定されていないと動作しないよー
  * 必要なレプリカ数=ceil(sum(podの現在のCPU使用率)/targetAverageUtilization)
  * スケールアウトの条件式: avg(Podの現在のCPU使用率)/targetAverageUtilization > 1.1
  * スケールインの条件式: avg(Podの現在のCPU使用率)/targetAverageUtilization < 0.9
  * 利用可能なメトリクスは下記
    * Resource: CPU/memoryのリソース
    * Object: k8s objectのメトリクス(ex. Ingressのhits/s)
    * Pods: Podのメトリクス(ex. Podのコネクション数)
    * External: k8s外のメトリクス(ex. LBのQPSなど)
* VerticalPodAutoscaler(VPA)はコンテナに割り当てるCPU/memoryのリソースの割り当てを自動的にスケールさせる