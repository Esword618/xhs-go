package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"net/url"
	"strconv"
	"time"
	"unsafe"
)

func Sign(uri string, data map[string]interface{}, ctime int64, a1, b1 string) map[string]string {
	v := strconv.FormatInt(ctime, 10)
	raw_str := fmt.Sprintf("%s%s%s", v, "test", uri)
	if dataJSON, err := json.Marshal(data); err == nil {
		raw_str += string(dataJSON)
	}

	md5Str := fmt.Sprintf("%x", md5.Sum([]byte(raw_str)))
	x_s := h(md5Str)
	x_t := v

	common := map[string]interface{}{
		"s0":  5, // getPlatformCode
		"s1":  "",
		"x0":  "1",     // localStorage.getItem("b1b1")
		"x1":  "3.2.0", // version
		"x2":  "Windows",
		"x3":  "xhs-pc-web",
		"x4":  "2.3.1",
		"x5":  a1, // cookie of a1
		"x6":  x_t,
		"x7":  x_s,
		"x8":  b1,             // localStorage.getItem("b1")
		"x9":  mrc(x_t + x_s), // Assuming mrc is implemented elsewhere
		"x10": 1,              // getSigCount
	}

	encodeStr, _ := json.Marshal(common)
	x_s_common := base64.StdEncoding.EncodeToString(encodeStr)

	return map[string]string{
		"x-s":        x_s,
		"x-t":        x_t,
		"x-s-common": x_s_common,
	}
}

//func Sign(a1, b1, x_s, x_t string) map[string]string {
//	common := map[string]any{
//		"s0":  5, // getPlatformCode
//		"s1":  "",
//		"x0":  "1",     // localStorage.getItem("b1b1")
//		"x1":  "3.3.0", // version
//		"x2":  "Windows",
//		"x3":  "xhs-pc-web",
//		"x4":  "1.4.4",
//		"x5":  a1, // cookie of a1
//		"x6":  x_t,
//		"x7":  x_s,
//		"x8":  b1, // localStorage.getItem("b1")
//		"x9":  mrc(x_t + x_s + b1),
//		"x10": 1, // getSigCount
//	}
//
//	encodeStr, _ := json.Marshal(common)
//	x_s_common := base64.StdEncoding.EncodeToString(encodeStr)
//	x_b3_traceid := GetB3TraceID()
//
//	return map[string]string{
//		"x-s":          x_s,
//		"x-t":          x_t,
//		"x-s-common":   x_s_common,
//		"x-b3-traceid": x_b3_traceid,
//	}
//}

func h(n string) string {
	m := ""
	d := "A4NjFqYu5wPHsO0XTdDgMa2r1ZQocVte9UJBvk6/7=yRnhISGKblCWi+LpfE8xzm3"
	for i := 0; i < 32; i += 3 {
		o := int(n[i])
		g := 0
		h := 0
		if i+1 < 32 {
			g = int(n[i+1])
		}
		if i+2 < 32 {
			h = int(n[i+2])
		}
		x := ((o & 3) << 4) | (g >> 4)
		p := ((15 & g) << 2) | (h >> 6)
		v := o >> 2
		b := h & 63
		if h == 0 {
			b = 64
		}
		if g == 0 {
			p = 64
			b = 64
		}
		m += string(d[v]) + string(d[x]) + string(d[p]) + string(d[b])
	}
	return m
}

func GetB3TraceID() string {
	re := "abcdef0123456789"
	je := 16
	e := ""

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 16; t++ {
		index := random.Intn(je)
		e += string(re[index])
	}

	return e
}

func mrc(e string) int32 {
	ie := []int{
		0, 1996959894, 3993919788, 2567524794, 124634137, 1886057615, 3915621685,
		2657392035, 249268274, 2044508324, 3772115230, 2547177864, 162941995,
		2125561021, 3887607047, 2428444049, 498536548, 1789927666, 4089016648,
		2227061214, 450548861, 1843258603, 4107580753, 2211677639, 325883990,
		1684777152, 4251122042, 2321926636, 335633487, 1661365465, 4195302755,
		2366115317, 997073096, 1281953886, 3579855332, 2724688242, 1006888145,
		1258607687, 3524101629, 2768942443, 901097722, 1119000684, 3686517206,
		2898065728, 853044451, 1172266101, 3705015759, 2882616665, 651767980,
		1373503546, 3369554304, 3218104598, 565507253, 1454621731, 3485111705,
		3099436303, 671266974, 1594198024, 3322730930, 2970347812, 795835527,
		1483230225, 3244367275, 3060149565, 1994146192, 31158534, 2563907772,
		4023717930, 1907459465, 112637215, 2680153253, 3904427059, 2013776290,
		251722036, 2517215374, 3775830040, 2137656763, 141376813, 2439277719,
		3865271297, 1802195444, 476864866, 2238001368, 4066508878, 1812370925,
		453092731, 2181625025, 4111451223, 1706088902, 314042704, 2344532202,
		4240017532, 1658658271, 366619977, 2362670323, 4224994405, 1303535960,
		984961486, 2747007092, 3569037538, 1256170817, 1037604311, 2765210733,
		3554079995, 1131014506, 879679996, 2909243462, 3663771856, 1141124467,
		855842277, 2852801631, 3708648649, 1342533948, 654459306, 3188396048,
		3373015174, 1466479909, 544179635, 3110523913, 3462522015, 1591671054,
		702138776, 2966460450, 3352799412, 1504918807, 783551873, 3082640443,
		3233442989, 3988292384, 2596254646, 62317068, 1957810842, 3939845945,
		2647816111, 81470997, 1943803523, 3814918930, 2489596804, 225274430,
		2053790376, 3826175755, 2466906013, 167816743, 2097651377, 4027552580,
		2265490386, 503444072, 1762050814, 4150417245, 2154129355, 426522225,
		1852507879, 4275313526, 2312317920, 282753626, 1742555852, 4189708143,
		2394877945, 397917763, 1622183637, 3604390888, 2714866558, 953729732,
		1340076626, 3518719985, 2797360999, 1068828381, 1219638859, 3624741850,
		2936675148, 906185462, 1090812512, 3747672003, 2825379669, 829329135,
		1181335161, 3412177804, 3160834842, 628085408, 1382605366, 3423369109,
		3138078467, 570562233, 1426400815, 3317316542, 2998733608, 733239954,
		1555261956, 3268935591, 3050360625, 752459403, 1541320221, 2607071920,
		3965973030, 1969922972, 40735498, 2617837225, 3943577151, 1913087877,
		83908371, 2512341634, 3803740692, 2075208622, 213261112, 2463272603,
		3855990285, 2094854071, 198958881, 2262029012, 4057260610, 1759359992,
		534414190, 2176718541, 4139329115, 1873836001, 414664567, 2282248934,
		4279200368, 1711684554, 285281116, 2405801727, 4167216745, 1634467795,
		376229701, 2685067896, 3608007406, 1308918612, 956543938, 2808555105,
		3495958263, 1231636301, 1047427035, 2932959818, 3654703836, 1088359270,
		936918000, 2847714899, 3736837829, 1202900863, 817233897, 3183342108,
		3401237130, 1404277552, 615818150, 3134207493, 3453421203, 1423857449,
		601450431, 3009837614, 3294710456, 1567103746, 711928724, 3020668471,
		3272380065, 1510334235, 755167117,
	}

	o := -1

	rightWithoutSign := func(num uint32, bit int) int {
		val := uint32(uintptr(unsafe.Pointer(&num))) >> bit
		MAX32INT := uint32(4294967295)
		return (int(val+(MAX32INT+1)) % (2 * int(MAX32INT+1))) - int(MAX32INT) - 1
	}

	for n := 0; n < 57; n++ {
		o = ie[(o&255)^int(e[n])] ^ rightWithoutSign(uint32(o), 8)
	}
	return int32(o ^ -1 ^ 3988292384)
}

var lookup = []string{
	"Z", "m", "s", "e", "r", "b", "B", "o", "H", "Q", "t", "N", "P", "+", "w", "O", "c",
	"z", "a", "/", "L", "p", "n", "g", "G", "8", "y", "J", "q", "4", "2", "K", "W", "Y",
	"j", "0", "D", "S", "f", "d", "i", "k", "x", "3", "V", "T", "1", "6", "I", "l", "U",
	"A", "F", "M", "9", "7", "h", "E", "C", "v", "u", "R", "X", "5",
}

func tripletToBase64(e uint32) string {
	return lookup[63&(e>>18)] +
		lookup[63&(e>>12)] +
		lookup[(e>>6)&63] +
		lookup[e&63]
}

func encodeChunk(e []byte, t, r int) string {
	m := make([]string, 0)
	for b := t; b < r; b += 3 {
		n := (16711680 & (int(e[b]) << 16)) +
			((int(e[b+1]) << 8) & 65280) +
			(int(e[b+2]) & 255)
		m = append(m, tripletToBase64(uint32(n)))
	}
	return fmt.Sprintf("%s", m)
}

func b64Encode(e string) string {
	P := len(e)
	W := P % 3
	U := make([]string, 0)
	z := 16383
	H := 0
	Z := P - W
	for H < Z {
		U = append(U, encodeChunk([]byte(e), H, min(Z, H+z)))
		H += z
	}
	if W == 1 {
		F := int(e[P-1])
		U = append(U, lookup[F>>2]+lookup[(F<<4)&63]+"==")
	} else if W == 2 {
		F := (int(e[P-2]) << 8) + int(e[P-1])
		U = append(U, lookup[F>>10]+lookup[63&(F>>4)]+lookup[(F<<2)&63]+"=")
	}
	return fmt.Sprintf("%s", U)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func encodeUtf8(e string) []int {
	b := []int{}
	m := url.QueryEscape(e)
	w := 0

	for w < len(m) {
		T := m[w]
		if T == '%' {
			E := m[w+1 : w+3]
			S := 0
			fmt.Sscanf(E, "%x", &S)
			b = append(b, S)
			w += 2
		} else {
			b = append(b, int(T))
		}
		w++
	}

	return b
}

func base36encode(number *big.Int, alphabet string) string {
	if number.Sign() < 0 {
		return "-" + base36encode(new(big.Int).Neg(number), alphabet)
	}

	base36 := ""
	length := big.NewInt(int64(len(alphabet)))

	zero := big.NewInt(0)
	div, mod := new(big.Int), new(big.Int)

	for number.Cmp(zero) != 0 {
		div.Mod(number, length)
		number.Div(number, length)
		mod.Abs(div)

		base36 = string(alphabet[mod.Int64()]) + base36
	}

	return base36
}

func base36decode(number string) *big.Int {
	result := new(big.Int)
	result.SetString(number, 36)
	return result
}

func GetSearchID() string {
	e := new(big.Int)
	timestamp := time.Now().UnixNano() / 1e6
	e.Lsh(big.NewInt(timestamp), 64)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	t := new(big.Int)
	t.SetInt64(int64(random.Intn(2147483646)))

	result := new(big.Int)
	result.Add(e, t)

	alphabet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return base36encode(result, alphabet)
}
