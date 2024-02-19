package main

import (
	"flag"
	"log"
	"net"
)

func main() {
	var ipArg string
	var instance string
	var port int

	flag.StringVar(&ipArg, "ip", "192.168.1.9", "ip of the printer")
	flag.IntVar(&port, "port", 1631, "port of the printer")
	flag.StringVar(&instance, "instance", "FPSTORESYD", "name of the printer")
	flag.Parse()

	var txt []string
	txt = append(txt, "txtvers=1")
	txt = append(txt, "qtotal=1")
	txt = append(txt, "pdl=application/octet-stream,image/urf,image/jpeg,image/pwg-raster")
	txt = append(txt, "rp=ipp/print")
	txt = append(txt, "note=")
	txt = append(txt, "ty=Brother HL-L3230CDW series")
	txt = append(txt, "product=(Brother HL-L3230CDW series)")
	txt = append(txt, "priority=25")
	txt = append(txt, "usb_MFG=Brother")
	txt = append(txt, "usb_MDL=HL-L3230CDW series")
	txt = append(txt, "usb_CMD=PJL,PCL,PCLXL,URF")
	txt = append(txt, "Color=T")
	txt = append(txt, "Copies=T")
	txt = append(txt, "Duplex=T")
	txt = append(txt, "Fax=F")
	txt = append(txt, "Scan=F")
	txt = append(txt, "PaperCustom=T")
	txt = append(txt, "Binary=T")
	txt = append(txt, "Transparent=T")
	txt = append(txt, "TBCP=F")
	txt = append(txt, "URF=SRGB24,W8,CP1,IS4-1,MT1-3-4-5-8-11,OB10,PQ4,RS600,V1.4,DM1")
	txt = append(txt, "kind=document,envelope,label,postcard")
	txt = append(txt, "PaperMax=legal-A4")
	txt = append(txt, "UUID=dbbd227c-d891-4474-9c92-7f20dcb234f4")
	txt = append(txt, "print_wfds=T")
	txt = append(txt, "mopria-certified=1.3")

	var ips []net.IP
	ip := net.ParseIP(ipArg)
	ips = append(ips, ip)

	service, err := NewMDNSService(instance, "_ipp._tcp", "", "", port, ips, txt)
	if err != nil {
		log.Fatalln(err)
	}

	// Create the mDNS server, defer shutdown
	server, _ := NewServer(&Config{
		Zone: service,
	})
	defer server.Shutdown()

	select {}
}
