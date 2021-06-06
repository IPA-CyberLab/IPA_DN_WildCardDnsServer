// Forked From: https://github.com/cunnie/sslip.io/
// (Golang-based DNS server which maps DNS records with embedded IP addresses to those addresses)
// by Brian Cunnie (https://github.com/cunnie/)

package main

import (
	"xip/xip"
)

func init_config() xip.Config {
	var cfg xip.Config

	// ドメイン名の指定
	// たとえば、 "12-34-56-78.abc.example.org" というドメイン名に対して 12.34.56.78 という IP アドレスを応答する
	// DNS システムを構築したい場合は、ここでは "abc.example.org." と表記する。
	// ドメイン名はすべての文字を小文字で記載する。
	// 文字列の最後は "." で終わる必要がある。
	cfg.DomainFqdn = "example.org."
	
	// SOA レコードに含まれるメールアドレス (DNS FQDN 形式の表記)。実質的に世界中で全く利用されていないので適当でよい。最後は "." で終わること。
	cfg.SoaMailAddress = "none.example.org."

	// ドメイン名に対する NS レコード。この wildcard_dns_server を動作させている権威 DNS サーバーのホスト名を指定する。
	// これは、上位ドメイン権威サーバー (例: レジストリの DNS サーバー) に登録されている NS レコードの設定と一致させること。
	// このサンプルでは、ns1 と ns2 の 2 つの NS が存在するものとして記述している。
	// 3 つ以上でも記述可能である。
	nsServerList := map[string]string {
		"ns1." + cfg.DomainFqdn: "1.2.3.4",
		"ns2." + cfg.DomainFqdn: "5.6.7.8",
	}

	// Let's Encrypt を用いたワイルドカード証明書更新サーバー IPA_DN_WildcardCertServerUtil (https://github.com/IPA-CyberLab/IPA_DN_WildcardCertServerUtil/) が動作している
	// サーバーの IP アドレスを指定する
	cfg.WildcardCertServerIp = "9.8.2.1"

	// このドメイン名そのものの A レコード、または "www." + このドメイン名そのものの A レコードの照会があったときに応答するアドレスを記述する。
	// すなわち、一般的なユーザーが Web ブラウザで http://ドメイン名/ にアクセスをしてみたときに Web ページを表示したい場合、
	// その Web サーバーのアドレスを記載するのである。
	cfg.DomainExactMatchARecord = "9.8.2.1"

	cfg.NsServerList = nsServerList
	return cfg
}

