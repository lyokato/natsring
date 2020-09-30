
## DESCRIPTION

"github.com/nats-io/nats.go"のクライアントのラッパ
複数のnats-serverに対して、hashringを使って分散してpubsubをかける

## EXAMPLE

より詳しくはexamples/main.goを参照のこと

```golang

  nodes := []string{
    "localhost:4222",
    "localhost:4223",
    "localhost:4224",
    "localhost:4226",
  }

  nr := natsring.New(nodes)

  // nats.ConnのConnectと違うのは、第一引数のURLは必要ということ
  // (Newで渡したすべてのノードに接続をかけるので)
  // あとは同じようにOptionを指定していく
  nr.ConnectAll(nats.MaxReconnects(10))


  // subscribe
  nr.Subscribe("MY_CHANNEL", func(m *nats.Msg) {
    fmt.Printf("RECEIVED: %s: %s\n", m.Subject, string(m.Data))
  })

  // publish
  err := nr.Publish("MY_CHANNEL", byte[]("FOOBAR"))


  // すべて切断
  nc.CloseAll()

```

Subscribe/Publishは、素のnats.Connとほぼ同じようにかけるが、
内部でHashRingを使い、チャンネルごとに違う接続を利用するようになっている。
