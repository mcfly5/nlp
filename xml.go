package main

import (
	//    "encoding/json"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	//"io/ioutil"
	"strconv"
	"strings"
)

type XMLDictionary struct {
	XMLName xml.Name `xml:"dictionary"`
	Lemmata Lemmata  `xml:"lemmata"`
	Links   Links    `xml:"links"`
}

type XMLDataStat struct {
	XMLName xml.Name `xml:"annotation"`
	L       []L      `xml:"text>paragraphs>paragraph>sentence>tokens>token>tfr>v>l"`
}

type Lemmata struct {
	XMLName  xml.Name `xml:"lemmata"`
	LemmList []Lemma  `xml:"lemma"`
}

type Links struct {
	XMLName  xml.Name `xml:"links"`
	LinkList []Link   `xml:"link"`
}

type Link struct {
	XMLName xml.Name `xml:"link"`
	From    string   `xml:"from,attr"`
	To      string   `xml:"to,attr"`
}

type Lemma struct {
	XMLName xml.Name `xml:"lemma"`
	Id      string   `xml:"id,attr"`
	Rev     string   `xml:"rev,attr"`
	L       L        `xml:"l"`
	F       []F      `xml:"f"`
}

type LemmaS struct {
	XMLName xml.Name `xml:"annotation"`
	L       []L      `xml:"text>paragraphs>paragraph>sentence>tokens>token>tfr>v>l"`
}

type L struct {
	T string `xml:"t,attr"`
	G []G    `xml:"g"`
}

type G struct {
	V string `xml:"v,attr"`
}

type F struct {
	T string `xml:"t,attr"`
}

type FirstForm struct {
	Form  string
	Gramm string
}

var dictForms map[string]map[int]bool
var dictLemmes map[int]FirstForm
var dictLinks map[int]int
var statDict map[string]map[string]int

func main() {

	//Reading and parssing XML from file
	//rawXmlData := "<data><person><firstname>Nic</firstname><lastname>Raboy</lastname><address><city>San Francisco</city><state>CA</state></address></person><person><firstname>Maria</firstname><lastname>Raboy</lastname></person></data>"
	//rawXmlData, _ := strconv.Unquote(setData())
	rawXmlData := readStringFromFile("dict.opcorpora.xml")
	//fmt.Printf(rawXmlData)
	var data XMLDictionary
	xml.Unmarshal([]byte(*rawXmlData), &data)

	fmt.Println("Lemmes in XML  :", len(data.Lemmata.LemmList))
	fmt.Println("Links  in XML  :", len(data.Links.LinkList))

	//Building the dictionary
	dictLemmes = make(map[int]FirstForm)
	dictForms = make(map[string]map[int]bool)
	dictLinks = make(map[int]int)

	var lemmForInsert FirstForm

	//Links
	var countBadLinks int
	for _, link := range data.Links.LinkList {

		from, err := strconv.Atoi(link.From)
		if err != nil {
			fmt.Println(err)
		}

		to, err := strconv.Atoi(link.To)
		if err != nil {
			fmt.Println(err)
		}

		if _, exists := dictLinks[to]; exists {
			countBadLinks++
		} else {
			dictLinks[to] = from
		}

	}

	fmt.Println("Bad links :", countBadLinks)

	//Lemmes
	for _, lemma := range data.Lemmata.LemmList {
		id, err := strconv.Atoi(lemma.Id)
		if err != nil {
			fmt.Println(err)
		}
		lemmForInsert.Form = lemma.L.T
		lemmForInsert.Gramm = lemma.L.G[0].V
		dictLemmes[id] = lemmForInsert
		//fmt.Println(lemma.L.T, " - ", lemma.L.G[0].V)
		for _, val := range lemma.F {
			//fmt.Println(val.T, " - ", lemma.L.T)
			if m, exist := dictForms[val.T]; exist {
				if _, exist := m[id]; !exist {
					dictForms[val.T][id] = true
				}
			} else {
				dictForms[val.T] = make(map[int]bool)
				dictForms[val.T][id] = true
			}
		}

	}

	fmt.Println("Lemmes in dict :", len(dictLemmes))
	fmt.Println("Forms in dict  :", len(dictForms))

	//Frequency
	rawXmlData = readStringFromFile("annot.opcorpora.no_ambig.xml")
	var dataStat XMLDataStat
	err := xml.Unmarshal([]byte(*rawXmlData), &dataStat)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	statDict = make(map[string]map[string]int)

	fmt.Println("Tokens:", len(dataStat.L))
	for _, k := range dataStat.L {
		if _, exists := statDict[k.T]; !exists {
			statDict[k.T] = make(map[string]int)
			statDict[k.T][k.G[0].V] = 1
		} else {
			statDict[k.T][k.G[0].V]++
		}
	}

	//Main
	strToSplit := `Стала стабильнее экономическая и политическая обстановка, предприятия вывели из тени зарплаты сотрудников.
	Все Гришины одноклассники уже побывали за границей, он был чуть ли не единственным, кого не вывозили никуда дальше Красной Пахры. Воркалось.`

	// Preparing the text
	r := strings.NewReplacer(".", " ", ",", " ", ";", " ", "!", " ", "?", " ")
	strToSplit = r.Replace(strToSplit)
	words := strings.Fields(strToSplit)
	fmt.Println(words)
	//fmt.Println(dictForms)

	for _, wordOrig := range words {
		word := strings.ToLower(wordOrig)
		str33 := ""
		maxStat := -1
		fmt.Println("\n", word, ":")
		for key, j := range dictForms[word] {

			fmt.Println("\t", key, j, ":")

			//str33 = str33 + ", " + strconv.Itoa(key) + dictLemmes[key].Form + "=" + dictLemmes[key].Gramm
			if val, exists := dictLinks[key]; exists {
				//str33 = str33 + "->" + strconv.Itoa(val) + dictLemmes[val].Form + "=" + dictLemmes[val].Gramm
				fmt.Print("\tclean:")
				str33 = str33 + "=" + dictLemmes[val].Gramm + "#"
				//				str33 = str33 + dictLemmes[val].Form + "=" + dictLemmes[val].Gramm
			} else {
				fmt.Print("\tdirt:")
				if maxStat < statDict[dictLemmes[key].Form][dictLemmes[key].Gramm] {
					str33 = str33 + "=" + dictLemmes[key].Gramm + "#" //+ ":" + strconv.Itoa(statDict[dictLemmes[key].Form][dictLemmes[key].Gramm])
					//str33 = str33 + dictLemmes[key].Form + "=" + dictLemmes[key].Gramm //+ ":" + strconv.Itoa(statDict[dictLemmes[key].Form][dictLemmes[key].Gramm])
					maxStat = statDict[dictLemmes[key].Form][dictLemmes[key].Gramm]
				}
			}
			fmt.Println("\t\t", str33)

		}
		//str33 = str33 + " ; Max:" + strconv.Itoa(maxStat)
		if len(dictForms[word]) == 0 {
			//Unknown word, need processing through trie
			str33 = "NI"
		}
		str33 = strings.Trim(str33, ",")
		str33 = strings.TrimSpace(str33)
		replacerGramm := strings.NewReplacer("NOUN", "S", "INFN", "V", "ADJF", "A", "PREP", "PR", "PRCL", "ADV", "ADVB", "ADV", "NPRO", "NI")
		str33 = replacerGramm.Replace(str33)

		//	fmt.Print(wordOrig, "{", str33, "} ")
		//		fmt.Print(word, " {", dictLemmes[dictForms[word][0]].Form, ", ", dictForms[word][0], ", ", dictLemmes[dictForms[word][0]].Gramm, "} ")

	}

	//jsonData, _ := json.Marshal(data)
	//fmt.Println(string(jsonData))
}

func readStringFromFile(file string) *string {

	//	bs2, err := ioutil.ReadFile("dict_test.xml")

	str := new(string)

	bs2, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error while opening a file #2...")
		return nil
	}
	*str = string(bs2)

	//println("Content #2:", str)
	/*

		str = strconv.Quote(`
			<dictionary>
				<lemmata>
					<lemma id="53399" rev="53399"><l t="всекитайский"><g v="ADJF"/></l><f t="всекитайский"><g v="masc"/><g v="sing"/><g v="nomn"/></f><f t="всекитайского"><g v="masc"/><g v="sing"/><g v="gent"/></f><f t="всекитайскому"><g v="masc"/><g v="sing"/><g v="datv"/></f><f t="всекитайского"><g v="anim"/><g v="masc"/><g v="sing"/><g v="accs"/></f><f t="всекитайский"><g v="inan"/><g v="masc"/><g v="sing"/><g v="accs"/></f><f t="всекитайским"><g v="masc"/><g v="sing"/><g v="ablt"/></f><f t="всекитайском"><g v="masc"/><g v="sing"/><g v="loct"/></f><f t="всекитайская"><g v="femn"/><g v="sing"/><g v="nomn"/></f><f t="всекитайской"><g v="femn"/><g v="sing"/><g v="gent"/></f><f t="всекитайской"><g v="femn"/><g v="sing"/><g v="datv"/></f><f t="всекитайскую"><g v="femn"/><g v="sing"/><g v="accs"/></f><f t="всекитайской"><g v="femn"/><g v="sing"/><g v="ablt"/></f><f t="всекитайскою"><g v="femn"/><g v="sing"/><g v="ablt"/><g v="V-oy"/></f><f t="всекитайской"><g v="femn"/><g v="sing"/><g v="loct"/></f><f t="всекитайское"><g v="neut"/><g v="sing"/><g v="nomn"/></f><f t="всекитайского"><g v="neut"/><g v="sing"/><g v="gent"/></f><f t="всекитайскому"><g v="neut"/><g v="sing"/><g v="datv"/></f><f t="всекитайское"><g v="neut"/><g v="sing"/><g v="accs"/></f><f t="всекитайским"><g v="neut"/><g v="sing"/><g v="ablt"/></f><f t="всекитайском"><g v="neut"/><g v="sing"/><g v="loct"/></f><f t="всекитайские"><g v="plur"/><g v="nomn"/></f><f t="всекитайских"><g v="plur"/><g v="gent"/></f><f t="всекитайским"><g v="plur"/><g v="datv"/></f><f t="всекитайских"><g v="anim"/><g v="plur"/><g v="accs"/></f><f t="всекитайские"><g v="inan"/><g v="plur"/><g v="accs"/></f><f t="всекитайскими"><g v="plur"/><g v="ablt"/></f><f t="всекитайских"><g v="plur"/><g v="loct"/></f></lemma>
					<lemma id="327807" rev="327807"><l t="собрание"><g v="NOUN"/><g v="inan"/><g v="neut"/></l><f t="собрание"><g v="sing"/><g v="nomn"/></f><f t="собранье"><g v="sing"/><g v="nomn"/><g v="V-be"/></f><f t="собрания"><g v="sing"/><g v="gent"/></f><f t="собранья"><g v="sing"/><g v="gent"/><g v="V-be"/></f><f t="собранию"><g v="sing"/><g v="datv"/></f><f t="собранью"><g v="sing"/><g v="datv"/><g v="V-be"/></f><f t="собрание"><g v="sing"/><g v="accs"/></f><f t="собранье"><g v="sing"/><g v="accs"/><g v="V-be"/></f><f t="собранием"><g v="sing"/><g v="ablt"/></f><f t="собраньем"><g v="sing"/><g v="ablt"/><g v="V-be"/></f><f t="собрании"><g v="sing"/><g v="loct"/></f><f t="собранье"><g v="sing"/><g v="loct"/><g v="V-be"/></f><f t="собраньи"><g v="sing"/><g v="loct"/><g v="V-be"/><g v="V-bi"/></f><f t="собрания"><g v="plur"/><g v="nomn"/></f><f t="собранья"><g v="plur"/><g v="nomn"/><g v="V-be"/></f><f t="собраний"><g v="plur"/><g v="gent"/></f><f t="собраниям"><g v="plur"/><g v="datv"/></f><f t="собраньям"><g v="plur"/><g v="datv"/><g v="V-be"/></f><f t="собрания"><g v="plur"/><g v="accs"/></f><f t="собранья"><g v="plur"/><g v="accs"/><g v="V-be"/></f><f t="собраниями"><g v="plur"/><g v="ablt"/></f><f t="собраньями"><g v="plur"/><g v="ablt"/><g v="V-be"/></f><f t="собраниях"><g v="plur"/><g v="loct"/></f><f t="собраньях"><g v="plur"/><g v="loct"/><g v="V-be"/></f></lemma>
					<lemma id="3278070" rev="327807"><l t="собраниеcdc"><g v="NOUN"/><g v="inan"/><g v="neut"/></l><f t="собрание"><g v="sing"/><g v="nomn"/></f><f t="собранье"><g v="sing"/><g v="nomn"/><g v="V-be"/></f><f t="собрания"><g v="sing"/><g v="gent"/></f><f t="собранья"><g v="sing"/><g v="gent"/><g v="V-be"/></f><f t="собранию"><g v="sing"/><g v="datv"/></f><f t="собранью"><g v="sing"/><g v="datv"/><g v="V-be"/></f><f t="собрание"><g v="sing"/><g v="accs"/></f><f t="собранье"><g v="sing"/><g v="accs"/><g v="V-be"/></f><f t="собранием"><g v="sing"/><g v="ablt"/></f><f t="собраньем"><g v="sing"/><g v="ablt"/><g v="V-be"/></f><f t="собрании"><g v="sing"/><g v="loct"/></f><f t="собранье"><g v="sing"/><g v="loct"/><g v="V-be"/></f><f t="собраньи"><g v="sing"/><g v="loct"/><g v="V-be"/><g v="V-bi"/></f><f t="собрания"><g v="plur"/><g v="nomn"/></f><f t="собранья"><g v="plur"/><g v="nomn"/><g v="V-be"/></f><f t="собраний"><g v="plur"/><g v="gent"/></f><f t="собраниям"><g v="plur"/><g v="datv"/></f><f t="собраньям"><g v="plur"/><g v="datv"/><g v="V-be"/></f><f t="собрания"><g v="plur"/><g v="accs"/></f><f t="собранья"><g v="plur"/><g v="accs"/><g v="V-be"/></f><f t="собраниями"><g v="plur"/><g v="ablt"/></f><f t="собраньями"><g v="plur"/><g v="ablt"/><g v="V-be"/></f><f t="собраниях"><g v="plur"/><g v="loct"/></f><f t="собраньях"><g v="plur"/><g v="loct"/><g v="V-be"/></f></lemma>
					<lemma id="174724" rev="174724"><l t="народный"><g v="ADJF"/><g v="Qual"/></l><f t="народный"><g v="masc"/><g v="sing"/><g v="nomn"/></f><f t="народного"><g v="masc"/><g v="sing"/><g v="gent"/></f><f t="народному"><g v="masc"/><g v="sing"/><g v="datv"/></f><f t="народного"><g v="anim"/><g v="masc"/><g v="sing"/><g v="accs"/></f><f t="народный"><g v="inan"/><g v="masc"/><g v="sing"/><g v="accs"/></f><f t="народным"><g v="masc"/><g v="sing"/><g v="ablt"/></f><f t="народном"><g v="masc"/><g v="sing"/><g v="loct"/></f><f t="народная"><g v="femn"/><g v="sing"/><g v="nomn"/></f><f t="народной"><g v="femn"/><g v="sing"/><g v="gent"/></f><f t="народной"><g v="femn"/><g v="sing"/><g v="datv"/></f><f t="народную"><g v="femn"/><g v="sing"/><g v="accs"/></f><f t="народной"><g v="femn"/><g v="sing"/><g v="ablt"/></f><f t="народною"><g v="femn"/><g v="sing"/><g v="ablt"/><g v="V-oy"/></f><f t="народной"><g v="femn"/><g v="sing"/><g v="loct"/></f><f t="народное"><g v="neut"/><g v="sing"/><g v="nomn"/></f><f t="народного"><g v="neut"/><g v="sing"/><g v="gent"/></f><f t="народному"><g v="neut"/><g v="sing"/><g v="datv"/></f><f t="народное"><g v="neut"/><g v="sing"/><g v="accs"/></f><f t="народным"><g v="neut"/><g v="sing"/><g v="ablt"/></f><f t="народном"><g v="neut"/><g v="sing"/><g v="loct"/></f><f t="народные"><g v="plur"/><g v="nomn"/></f><f t="народных"><g v="plur"/><g v="gent"/></f><f t="народным"><g v="plur"/><g v="datv"/></f><f t="народных"><g v="anim"/><g v="plur"/><g v="accs"/></f><f t="народные"><g v="inan"/><g v="plur"/><g v="accs"/></f><f t="народными"><g v="plur"/><g v="ablt"/></f><f t="народных"><g v="plur"/><g v="loct"/></f></lemma>
					<lemma id="267083" rev="267083"><l t="представитель"><g v="NOUN"/><g v="anim"/><g v="masc"/></l><f t="представитель"><g v="sing"/><g v="nomn"/></f><f t="представителя"><g v="sing"/><g v="gent"/></f><f t="представителю"><g v="sing"/><g v="datv"/></f><f t="представителя"><g v="sing"/><g v="accs"/></f><f t="представителем"><g v="sing"/><g v="ablt"/></f><f t="представителе"><g v="sing"/><g v="loct"/></f><f t="представители"><g v="plur"/><g v="nomn"/></f><f t="представителей"><g v="plur"/><g v="gent"/></f><f t="представителям"><g v="plur"/><g v="datv"/></f><f t="представителей"><g v="plur"/><g v="accs"/></f><f t="представителями"><g v="plur"/><g v="ablt"/></f><f t="представителях"><g v="plur"/><g v="loct"/></f></lemma>
					<lemma id="155465" rev="155465"><l t="март"><g v="NOUN"/><g v="inan"/><g v="masc"/></l><f t="март"><g v="sing"/><g v="nomn"/></f><f t="марта"><g v="sing"/><g v="gent"/></f><f t="марту"><g v="sing"/><g v="datv"/></f><f t="март"><g v="sing"/><g v="accs"/></f><f t="мартом"><g v="sing"/><g v="ablt"/></f><f t="марте"><g v="sing"/><g v="loct"/></f><f t="марты"><g v="plur"/><g v="nomn"/></f><f t="мартов"><g v="plur"/><g v="gent"/></f><f t="мартам"><g v="plur"/><g v="datv"/></f><f t="марты"><g v="plur"/><g v="accs"/></f><f t="мартами"><g v="plur"/><g v="ablt"/></f><f t="мартах"><g v="plur"/><g v="loct"/></f></lemma>

				</lemmata>
			</dictionary>
			`)
		//*/
	return str
}
