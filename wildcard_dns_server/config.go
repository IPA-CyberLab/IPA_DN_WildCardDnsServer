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
	cfg.DomainFqdn = "test.com."
	
	// SOA レコードに含まれるメールアドレス (DNS FQDN 形式の表記)。通常は変更する必要はない。最後は "." で終わること。
	cfg.SoaMailAddress = "postmaster.example.jp."

	// ドメイン名に対する NS レコード。この wildcard_dns_server を動作させている権威 DNS サーバーのホスト名を指定する。
	// このサンプルでは、ns1 と ns2 の 2 つの NS が存在するものとして記述している。
	// 3 つ以上でも記述可能である。
	nsServerList := map[string]string {
		"ns1." + cfg.DomainFqdn: "1.2.3.1",
		"ns2." + cfg.DomainFqdn: "1.2.3.2",
		"ns3." + cfg.DomainFqdn: "1.2.3.3",
	}

	cfg.NsServerList = nsServerList
	return cfg
}

