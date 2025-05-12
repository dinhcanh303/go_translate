package go_translate

var GoogleUrls map[GoogleAPIType]string = map[GoogleAPIType]string{
	TypeHtml:               "https://translate-pa.googleapis.com/v1/translateHtml",
	TypeClientDictChromeEx: "/translate_a/t?client=dict-chrome-ex&sl=auto",
	TypeClientGtx:          "/translate_a/single?client=gtx&sl=auto&dt=t",
	TypePaGtx:              "https://translate-pa.googleapis.com/v1/translate?params.client=gtx&query.source_language=auto&data_types=TRANSLATION&data_types=SENTENCE_SPLITS&data_types=BILINGUAL_DICTIONARY_FULL",
}

const MicrosoftServerUrl = "https://webmail.smartlinkcorp.com/dotrans_20160909.php"

var GoogleAPITypeSupport = []GoogleAPIType{TypeHtml, TypeClientGtx, TypeClientDictChromeEx, TypePaGtx}

var MpGoogleAPITypeSupport = map[GoogleAPIType]struct{}{TypeHtml: {}, TypeClientGtx: {}, TypeClientDictChromeEx: {}, TypePaGtx: {}, TypeRandom: {}, TypeSequential: {}, TypeMix: {}}

var DefaultServiceUrls = []string{
	"translate.google.com",
	"translate.google.mk",     //152ms-50ms
	"translate.google.co.th",  //167ms-55ms
	"translate.google.com.co", //63ms-211ms
	"translate.google.com.mx", //57ms-149ms
	"translate.google.com.tj", //189ms-51ms
	"translate.google.cv",     //186ms-95ms
	"translate.google.gp",     //185ms-74ms
	"translate.google.ki",     //175ms-70ms
	"translate.google.ps",     //174ms-68ms
	"translate.google.com.py", //167ms-52ms
	"translate.google.com.gh", //167ms-53ms
	"translate.google.com.gi", //171ms-110ms
	"translate.google.com.uy", //177ms-53ms
	"translate.google.com.vc", //146ms-49ms
	"translate.google.hu",     //171ms-59ms
	"translate.google.bi",     //181ms-59ms
	"translate.google.dz",     //193ms-65ms
	"translate.google.com.bz", //194ms-52ms
	"translate.google.com.sa", //194ms-57ms
	"translate.google.kg",     //195ms-53ms
	"translate.google.cm",     //172ms-178ms
	"translate.google.co.id",  //180ms-73ms
	"translate.google.co.tz",  //195ms-62ms
	"translate.google.com.cu", //178ms-82ms
	// "translate.google.com.ly", //197ms-74ms
	// "translate.google.je",     //194ms-227ms
	// "translate.google.com.gt", //164ms-227ms
	// "translate.google.com.sl", //116ms-218ms
	// "translate.google.com.vn", //60ms-208ms
	// "translate.google.ad",     //217ms-232ms
	// "translate.google.be",     //217ms-208ms
	// "translate.google.co.cr",  //212ms-116ms
	// "translate.google.co.mz",  //224ms-81ms
	// "translate.google.co.vi",  //279ms-68ms
	// "translate.google.com.ai", //202ms-54ms
	// "translate.google.com.cy", //286ms-55ms
	// "translate.google.com.lb", //225ms-70ms
	// "translate.google.com.my", //206ms-59ms
	// "translate.google.com.pk", //272ms-177ms
	// "translate.google.cz",     //227ms-73ms
	// "translate.google.lv",     //202ms-55ms
	// "translate.google.ne",     //205ms-70ms
	// "translate.google.rw",     //216ms-71ms
	// "translate.google.se",     //233ms-84ms
	// "translate.google.sh",     //248ms-83ms
	// "translate.google.si",     //255ms-63ms
	// "translate.google.sk",     //242ms-100ms
	// "translate.google.co.zw",  //201ms-108ms
	// "translate.google.co.nz",  //230ms-61ms
	// "translate.google.com.hk", //201ms-50ms
	// "translate.google.tn",     //209ms-66ms
	// "translate.google.tt",     //237ms-62ms
	// "translate.google.at",     //5.17s-223ms
}

var DefaultUserAgents = []string{
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/123.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_2) AppleWebKit/605.1.15 Version/16.0 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 Chrome/105.0.0.0 Safari/537.36 Edge/18.19041",
	"Mozilla/5.0 (Android; Mobile; rv:27.0) Gecko/27.0 Firefox/27.0",
	"Mozilla/5.0 (Linux; Android 9; SM-T820 Build/PPR1.180610.011; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/85.0.4183.101 Safari/537.36 ANDROID_APP",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Silk/85.3.5 like Chrome/85.0.4183.126 Safari/537.36",
	"BlackBerry9700/5.0.0.862 Profile/MIDP-2.1 Configuration/CLDC-1.1 VendorID/331 UNTRUSTED/1.0 3gpp-gba",
	"Mozilla/5.0 (BlackBerry; U; BlackBerry 9930; en-US) AppleWebKit/534.11+ (KHTML, like Gecko) Version/7.0.0.241 Mobile Safari/534.11+",
	"Mozilla/5.0 (BlackBerry; U; BlackBerry 9800; zh-TW) AppleWebKit/534.8+ (KHTML, like Gecko) Version/6.0.0.448 Mobile Safari/534.8+",
	"Mozilla/5.0 (Linux; U; Android 4.0.4; en-gb; GT-I9300 Build/IMM76D) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
	"Mozilla/5.0 (Linux; U; Android 2.2; en-us; SCH-I800 Build/FROYO) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
	"Mozilla/5.0 (Linux; Android 10.0; X116L) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36",
	"Mozilla/5.0 (X11; U; U; Linux x86_64; nl-nl) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36 Puffin/8.4.0.42081AP",
	"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; Trident/6.0)",
	"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; Trident/6.0; Touch)",
	"Mozilla/5.0(iPad; U; CPU iPhone OS 3_2 like Mac OS X; en-us) AppleWebKit/531.21.10 (KHTML, like Gecko) Version/4.0.4 Mobile/7B314 Safari/531.21.10",
	"Mozilla/5.0 (X11; U; Linux x86_64; nl-NL) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.114 Safari/537.36 Puffin/5.2.3IT",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 5_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9B179 Safari/7534.48.3",
	"Opera/9.80 (J2ME/MIDP; Opera Mini/9.80 (J2ME/22.478; U; en) Presto/2.5.25 Version/10.54",
	"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_3; en-us; Silk/1.1.0-80) AppleWebKit/533.16 (KHTML, like Gecko) Version/5.0 Safari/533.16 Silk-Accelerated=true",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_6_8) AppleWebKit/537.13+ (KHTML, like Gecko) Version/5.1.7 Safari/534.57.2",
	"Mozilla/5.0 (Linux; Android 4.4.4; Nexus 5 Build/KTU84P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.117 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 4.1.1; Nexus 7 Build/JRO03D) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.166  Safari/535.19",
	"Mozilla/5.0 (MeeGo; NokiaN9) AppleWebKit/534.13 (KHTML, like Gecko) NokiaBrowser/8.5.0 Mobile Safari/534.13",
	"Mozilla/5.0 (SymbianOS/9.4; Series60/5.0 NokiaN97-1/12.0.024; Profile/MIDP-2.1 Configuration/CLDC-1.1; en-us) AppleWebKit/525 (KHTML, like Gecko) BrowserNG/7.1.12344",
	"Mozilla/5.0 (X11; U; Linux armv7l; no-NO; rv:1.9.2.3pre) Gecko/20100723 Firefox/3.5 Maemo Browser 1.7.4.8 RX-51 N900",
	"Mozilla/5.0 (PlayBook; U; RIM Tablet OS 2.0.1; en-US) AppleWebKit/535.8+ (KHTML, like Gecko) Version/7.2.0.1 Safari/535.8+",
	"Mozilla/5.0 (PLAYSTATION 3 4.60) AppleWebKit/531.22.8 (KHTML, like Gecko)",
	"Mozilla/5.0 (PlayStation Vita 3.12) AppleWebKit/536.26 (KHTML, like Gecko) Silk/3.2",
	"Mozilla/5.0 (Linux; Android 4.4.2; en-us; SAMSUNG SCH-I545 Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Version/1.5 Chrome/28.0.1500.94 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 4.1.2; GT-I8190 Build/JZO54K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.117 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 4.4.2; en-gb; SAMSUNG SM-G900F Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Version/1.6 Chrome/28.0.1500.94 Mobile Safari/537.36",
	"Mozilla/5.0 (SAMSUNG; SAMSUNG-GT-S8530/S8530DDLC2; U; Bada/2.0; en-us) AppleWebKit/534.20 (KHTML, like Gecko) Dolfin/3.0 Mobile WVGA SMM-MMS/1.2.0 OPN-B",
	"Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:15.0) Gecko/20100101 Firefox/15.0.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/70.0.3538.64 Safari/537.36Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:83.0) Gecko/20100101 Firefox/83.0",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows Phone OS 7.5; Trident/5.0; IEMobile/9.0; SAMSUNG; SGH-i917)",
	"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; ARM; Trident/6.0)",
	"SAMSUNG-SGH-I617/UCHJ1 Mozilla/4.0 (compatible; MSIE 6.0; Windows CE; IEMobile 7.11)",
	"Mozilla/4.0 (compatible; MSIE 6.0; Windows CE; IEMobile 8.12; MSIEMobile 6.0) 320x240; VZW; UTStar-XV6175.1; Windows Mobile 6.5 Standard;",
	"Opera/9.80 (Android 2.3.3; Linux; Opera Mobi/ADR-1202011015; U; en) Presto/2.9.201 Version/11.50",
	"Opera/9.80 (BREW; Opera Mini/5.0/27.2370; U; en) Presto/2.8.119 240X320 Samsung SCH-U380",
	"Mozilla/5.0 (Windows; U; Win 9x 4.90; en-GB; rv:1.8.1.1) Gecko/20061204 Firefox/2.0.0.1",
	"Mozilla/5.0 (X11; U; SunOS sun4u; en-US; rv:1.6) Gecko/20040503",
	"Mozilla/5.0 (X11; CrOS i686 0.12.433) AppleWebKit/534.30 (KHTML, like Gecko) Chrome/12.0.742.77 Safari/534.30",
	//Google bot
	"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	"Mozilla/5.0 (iPhone; U; CPU iPhone OS) (compatible; Googlebot-Mobile/2.1; http://www.google.com/bot.html)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	"DoCoMo/2.0 N905i(c100;TB;W24H16) (compatible; Googlebot-Mobile/2.1; +http://www.google.com/bot.html)",
}

const GOOGLE_API_KEY_TRANSLATE = "AIzaSyATBXajvzQLTDHEQbcpq0Ihe0vWDHmO520"
const GOOGLE_API_KEY_TRANSLATE_PA = "AIzaSyDLEeFI5OtFBwYBIoK_jj5m32rZK5CkCXA"
