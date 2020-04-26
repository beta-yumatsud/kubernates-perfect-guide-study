## 章立て
コンテナのログ収集

## memo
* 前提として、コンテナのアプリケーションログは標準出力、標準エラー出力にログを出力しておくこと

### Fluentd
* DaemonSetを利用して各ノードにFluentdのPodを1つずつ配置しておけばOK
* 転送先は色々と選べるよー
* fluent bitというC言語で書かれた軽量版fluentdもあるみたい（IoT用とかに開発されたらしい）

### Datadog
* Datadog Logsでもログ収集は可能
* こちらはDatadog Agentが収集と転送を行う
