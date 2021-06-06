# IPA_DN_WildCardDnsServer (Go 言語で書かれたステートレスなワイルドカード DNS サーバー)
作成: 登 大遊 (Daiyuu Nobori)


## 背景
インターネット上で動作するシステムを開発するとき、動的に増減するサーバーをグローバル IPv4 / IPv6 アドレス上で動作させたい場合で、かつ、これらの IPv4 / IPv6 アドレスがクラウドシステムやダイナミック IP アドレス割り当て方式の ISP などのように事前に固定できない場合がある。


ところが、最近の Web ブラウザに代表される色々なシステムは、IP アドレスを URL 部分に直接指定した形式の URI を扱いづらい。

たとえば、HTTP や HTTPS を用いるとき、システム内でのサーバーを指定するアドレスとして `https://1.2.3.4/` とか `wss://1.2.3.4/` のような指定を行なうと、Web ブラウザ等のセキュリティ機構によって色々な問題が発生することが多くある。

この問題を解決するための一般的に普及している手法は、動的にレコードを管理するダイナミック DNS システムや、リバースプロキシに代表されるようなロードバランサ等を用いる方法である。しかし、ダイナミック DNS システムは、状態 (データ) を保持する必用があり、メンテナンスや保守のコストがかかる。また、アドレス更新のために何らかのクレデンシャルを扱う必要があり、セキュリティ上の手間がかかる。リバースプロキシに代表されるようなロードバランサ等を用いる場合は、これらの中継システムに負荷がかかったり、遅延の発生原因となったりする。また、バックエンドのシステムとこれらの中継システムとの間の連携が必要になる。

そこで、上記のような追加の手間が不要な方法として、IPv4 アドレスまたは IPv6 アドレスが埋め込まれた FQDN を DNS 照会すると、埋め込み IPv4 アドレスまたは IPv6 アドレスが応答するようなシンプルな DNS サーバーがあれば便利である。この手法を用いたフリーのサービスとして、以前より http://xip.io/ や https://sslip.io/ がある。

sslip.io の実装である https://github.com/cunnie/sslip.io/ には、下記の既存の課題があった。

- 設定パラメータがソースコード中の色々な部分に埋め込まれており、変更が容易でない。
- Let's Encrypt のような ACME 証明書自動発行サーバーに対応するための _acme_challenge に対応する NS レコードとして、予め設定されているレコードを応答することができない。


## 目的
本プログラムの目的は、sslip.io のプログラムを一部改良し、以下の機能を満たす、任意の Linux サーバー上で稼働可能な DNS コンテンツ・サーバーを実装するものである。

以下の説明では、本プログラムを `example.org` というドメインで使用する場合を例示している。この `example.org` という文字列を、実際に運用したいドメイン名に置換して解釈すること。

- 本プログラムの派生元 (sslip.io) の仕組みを継承した機能
  - `example.org` などのある単一のドメイン (既存のドメインのサブドメインでもよい) を所有しているとき、たとえば、`1-2-3-4.example.org.` という FQDN に対する IPv4 (A) レコード照会への応答として `1.2.3.4` という IPv4 (A) レコードを返す。  
    このように、IPv4 アドレスの 4 バイトを 10 進数の数字で表記し、数字の値をハイフン `(-)` でつなぐことで、単一の IPv4 アドレスを表すのである。
    - また、`websocket-1-2-3-4.example.org` というように、IPv4 アドレス部分の前に任意の英数字のプレフィックスを挿入することも可能である。ここで、本 DNS サーバーはプレフィックス部分を無視し、この例においては `1.2.3.4` を返す。
  - IPv6 もサポートしている。たとえば、`2001:1:2:3::4` という IPv6 アドレスを返したい場合は、`2001-1-2-3--4` のように、本来の IPv6 アドレスのコロン部分をハイフンに置換し FQDN を生成すればよい。
  - その他の基本的な仕組みと規則は https://sslip.io/ を参照すること。
- 本プログラムの独自の機能
  - Let's Encrypt ワイルドカード証明書発行および提供サーバープログラムである IPA_DN_WildcardCertServerUtil (https://github.com/IPA-CyberLab/IPA_DN_WildcardCertServerUtil) との連携を可能にした。
  - Let's Ecnrypt の検証サーバーから `_acme_challenge.example.org` に対して TXT レコードの要求があった場合は、IPA_DN_WildcardCertServerUtil が動作しているサーバー (このサーバーは、本プログラムが動作しているサーバーとは別のサーバーである必要がある) のグローバル IPv4 アドレスを埋め込んだ `ssl-cert-server.example.org` という A レコードを指す NS レコードを応答する。これにより、Let's Encrypt によるワイルドカード証明書 `*.example.org` の発行申請とその検証が可能となる。
  - 上記の手法による Let's Encrypt ワイルドカード証明書発行との組み合わせにより、以下のことが可能となる。
    1. HTTPS WebSocket などの、自己署名証明書を利用することが不可能な局面において、`websocket-1-2-3-4.example.org` というような FQDN をステートレスに発行することができ、Web ブラウザのエラーを回避することができる。
    2. 上記 1. における WebSocket を提供するサーバーがクライアントに提示する SSL 証明書は、IPA_DN_WildcardCertServerUtil サーバーが Let's Encrypt を用いて 1 ヶ月に 1 回程度更新を試みた結果取得している最新版のワイルドカード SSL 証明書 `*.example.org` を利用することができる。ここで、この証明書の X.509 証明書本体 (.crt) および秘密鍵 (.key) ファイルは、証明書は、IPA_DN_WildcardCertServerUtil サーバー上の nginx で動作しているインチキ・Web ディレクトリ上に置かれている。これを各利用サーバーは Basic 認証 (SSL 経由でアクセスできるはずなので安全である) で取得してダウンロードし、利用することができる。この Basic 認証を適正に行なうことにより、ワイルドカード証明書の秘密鍵が漏えいするセキュリティ問題は防止できる。なお、万一漏えいをしてしまった場合であっても、Let's Encrypt の証明書有効期限は短いので、被害は限定的である。
    3. 上記のすべては、ほぼステートレス処理である。すなわち、本プログラムを動作させるサーバーと、IPA_DN_WildcardCertServerUtil を動作させるサーバーには、大した重要な情報は管理されない。これらのサーバーは、極めて安価なクラウドサーバー (例: Amazon EC2 の VM) 上で立ち上げることができる。また、万一クラウドサーバーがうまく動作しなくなった場合も、直ちに別のサーバーで立ち上げることができる。これにより、保守運用のコストが削減できる。
  - `example.org` そのもの、および `www.example.org` というアドレスに対する A レコード要求があったときは、設定された IPv4 アドレスを固定的に応答する。
  - すべての設定項目は、`config.go` という単一かつ単純なテキストファイルを直接編集することによって設定・変更が可能である。


# インストールマニュアル
## 必要なもの
事前に必要なものは、以下のとおりである。
- 本プログラムによって運用したいドメイン名 1 個  
  (上記の説明における `example.org` に相当するもの。サブドメインでもよい。)
- 本プログラムを動作させる実際の NS サーバー (権威サーバー) として、現代的な Linux が動作する任意のクラウドまたはオンプレミスの VM であって、固定グローバル IPv4 アドレスの割当てがされているもの 2 台以上  
  (原理的には 1 台でもよいが、冗長のために 2 台以上を推奨する。また、ドメインレジストリによっては、2 台以上の NS レコードを登録しなければならない規則のところがあり、このような場合は、必然的に 2 台以上のサーバーが必要となる。)
  - Linux のバージョンは、Ubuntu 20.04 または Ubuntu 18.04 を推奨する。それ以外の Linux でもおおむね動作すると思われるが、自己責任で動作させること。
- もし、Let's Encrypt ワイルドカード証明書発行および提供サーバープログラムである IPA_DN_WildcardCertServerUtil (https://github.com/IPA-CyberLab/IPA_DN_WildcardCertServerUtil) との連携をしたい場合は、これとは別に IPA_DN_WildcardCertServerUtil をインストールし設定した状態のサーバー 1 台。これは後からセットアップしてもよい。


## AWS インスタンスの作成 (AWS を利用する場合)
このマニュアルでは、2 台の VM として、以下を用意するものとして説明を行なう。
```
Amazon EC2 の最も安価な VM インスタンス 2 台。ARM64 でも x64 でもよい。
2021/06/06 時点では、「t4g.nano」インスタンス (ARM 64bit CPU, RAM 0.5GB) が最もランニングコストが安価である。
これは 0.0042 USD / 時間 (例: 米オレゴンまたは米バージニア北部の DC を選択した場合) のため、1 ドル 120 円として、1 ヶ月あたり 375 円で 1 台の VM を運用することができる。
一方、東京の DC の場合は、0.0054 USD / 時間のため、1 ヶ月あたり 483 円で 1 台の VM を運用することができる。
これらの 2 台のインスタンスは、冗長のため、できるだけ、別々のアベイビリティゾーンで稼働させることを推奨する。
上記の EC2 コストは、コンピューティングコストであり、ストレージやネットワークは別途課金が発生する可能性があるため注意すること。また、最新のコストは AWS の Web サイトを確認すること。
```


これらの 2 台の VM は、以下のような設定で作成する。
```
AMI: Ubuntu Server 20.04 LTS (HVM), SSD Volume Type - 64 ビット (Arm)
インスタンスタイプ: t4g.nano
ネットワークのセキュリティグループ: SSH (TCP 22)、DNS (UDP 53)、ICMP IPv4 (楽しみのため) のみを任意のソース IP から通す。
Elastic IP をそれぞれの VM 用に作成し、各 VM に固定で関連付ける。
```

## VM の設定 (SSH 経由)
2 台の VM それぞれについて、SSH 経由で以下のように設定する。


まず、作成したばかりの EC2 サーバーに、ユーザー `ubuntu` として SSH サーバーにログインする。

```
# タイムゾーンを日本標準時 (Asia/Tokyo) に設定する
sudo timedatectl set-timezone Asia/Tokyo

# 最近の Ubuntu はヘンなローカル DNS プロキシが動作しており、けしからん。
# これらを以下のように停止するのである。
sudo systemctl disable systemd-resolved
sudo systemctl stop systemd-resolved
sudo rm /etc/resolv.conf

# すると resolv.conf がなくなってしまうので、
# インチキ Google Public DNS サーバーをひとまず手動で設定する。
echo nameserver 8.8.8.8 | sudo tee /etc/resolv.conf

# 上記を設定すると、Linux が sudo するたびに
# sudo: unable to resolve host ip-xxx-xxx-xxx-xxx: Name or service not known
# などと言ってくるようになりうっとおしいので、
# /etc/hosts に自ホストを追記して解決する。
echo $(ip route get 8.8.8.8 | cut -d " " -f 3 | head -n 1) $(hostname) | sudo tee /etc/hosts

# apt-get でいやないやな go 言語をインストールする。
sudo apt-get -y update && sudo apt-get -y install golang-go

# git で 本 DNS サーバープログラム (IPA_DN_WildCardDnsServer) をダウンロードする。
sudo mkdir -p /opt/IPA_DN_WildCardDnsServer/
sudo chown ubuntu:ubuntu /opt/IPA_DN_WildCardDnsServer/
cd /opt/IPA_DN_WildCardDnsServer/
git clone https://github.com/IPA-CyberLab/IPA_DN_WildCardDnsServer.git

# テキストエディタで設定ファイルをいじる。コメントを参照すること。
nano /opt/IPA_DN_WildCardDnsServer/IPA_DN_WildCardDnsServer/wildcard_dns_server/config.go

# 本 DNS サーバープログラムをテスト実行をしてみる。
cd /opt/IPA_DN_WildCardDnsServer/IPA_DN_WildCardDnsServer/wildcard_dns_server/
sudo go run main.go config.go

# 糸冬了！！
```

## DNS サーバーが正しく動作していることのテスト
手元の Windows, Linux, mac 等で古くさい nslookup コマンドを用いて以下のようにテストする。
```
C:\>nslookup

> server <VM のグローバル IPv4 アドレス>

> set q=a
> 1-2-3-4.example.org.
名前:    1-2-3-4.example.org
Address:  1.2.3.4

> ahobakamanuke-1-2-3-4.example.org.
名前:    ahobakamanuke-1-2-3-4.example.org
Address:  1.2.3.4

> set q=aaaa
> 2001-1-2-3-4--9821.example.org.
名前:    2001-1-2-3-4--9821.example.org
Address:  2001:1:2:3:4::9821

> set q=ns
> example.org.
example.org     nameserver = ns1.example.org
example.org     nameserver = ns2.example.org
ns1.example.org internet address = 1.2.3.4
ns2.example.org internet address = 5.6.7.8

> set q=txt
> _acme-challenge.example.org.
_acme-challenge.example.org     nameserver = ssl-cert-server.example.org
ssl-cert-server.example.org     internet address = 9.8.2.1

```

問題なければ、テスト実行している go を Ctrl + C で終了する。


## デーモン化
IPA_DN_WildCardDnsServer プログラムが Linux 起動後に自動的にデーモンとして動作するようにするためには、以下のようにする。
```
# デーモン定義ファイルの作成
# (少し長いが、EOF の部分までコピーペーストして SSH で一気に貼り付けること。)

sudo dd of=/etc/systemd/system/IPA_DN_WildCardDnsServers.service <<\EOF
[Unit]
Description=IPA_DN_WildCardDnsServers
After=network.target
StartLimitIntervalSec=0

[Service]
Type=forking
User=root
Group=root
ExecStart=/usr/bin/go run main.go config.go
WorkingDirectory=/opt/IPA_DN_WildCardDnsServer/IPA_DN_WildCardDnsServer/wildcard_dns_server/
Restart=always
Type=simple
User=root
Environment=HOME=/root
StandardOutput=null
RestartSec=5s

[Install]
WantedBy=multi-user.target

EOF

# デーモン定義ファイルの読み込み
sudo systemctl daemon-reload

# システム起動時に自動起動するように設定
sudo systemctl enable IPA_DN_WildCardDnsServers

# 動作開始
sudo systemctl start IPA_DN_WildCardDnsServers

sudo systemctl stop IPA_DN_WildCardDnsServers

# 正しく動作しているかどうか確認
sudo systemctl status IPA_DN_WildCardDnsServers

# 糸冬了！！
```


デーモン化が完了したら、reboot して自動的に DNS サーバーが動作開始することを確認する。

また、上記の `nslookup` コマンドによる外部クライアント端末からの実際の DNS クエリテストも実行し、正しく動作することを確認する。


## 上位 DNS サーバー (レジストリ等の DNS サーバー) への NS 委譲レコードの登録申請
上記のように DNS サーバー (少なくとも 2 台) の稼働を開始させた後に、上流の DNS サーバーへの NS レコードの登録を申請する。

ここで、NS レコードの名前と IP アドレスは `config.go` で指定したものと一致させる必要があるため、注意すること。

上位 DNS サーバーへの登録が行なわれると、無事、インターネットから名前解決ができるようになる。

これにより、ワイルドカード DNS サーバーの稼働が開始される。

なお、設定ファイル `config.go` の内容を変更した場合は、デーモンの再起動 (面倒であれば VM の再起動) が必要である。そして、再起動後には、必ず、正しく動作しているかどうか十分よく確認すること。



