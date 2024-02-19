# AirPrinter

A fake AirPrint printer that reverse proxies into another AirPrint printer.
The reverse proxy modifies the IPP attributes so that iOS thinks the printer has A6 loaded in the paper tray.
Now you can print A6 from whenever and wherever you want.

All this because iOS gimped the print dialog and does not allow us a way to specify the paper size for printers that have a manual tray.

## Features
* Works from MacOS and iOS
* `mdns` resolver that responds to `_ipp._tcp` and `_universal._sub._ipp._tcp`
* Hardcoded TXT records for my printer and A6 paper size

## How to debug

Use `dns-sd -Z _ipp._tcp` and `dns-sd -Z _ipp._tcp,universal` to verify mDNS functionality, use it dump TXT records from real printers.
Use `rvictl` to check traffic logs between the fake printer and iOS.

## Notes

Because this runs its own mDNS, it is incompatible with Avahi or any other mDNS resolvers that you may be already running.
Make sure to turn off these services.

iOS insists that the IPP port has to be on 631 and will not accept any other ports.

## References

* [Bonjour AirPrint Spec](https://developer.apple.com/bonjour/printing-specification/bonjourprinting-1.2.1.pdf)
* [AirPrint MDM payload settings for Apple devices](https://support.apple.com/en-us/guide/deployment/dep3b4cf515/web)
