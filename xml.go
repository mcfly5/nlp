package main

import (
	//    "encoding/json"

	"encoding/xml"
	"fmt"
	"strings"
	//"io/ioutil"
	"strconv"
)

type Lemmata struct {
	XMLName  xml.Name `xml:"lemmata"`
	LemmList []Lemma  `xml:"lemma"`
}

type Lemma struct {
	XMLName xml.Name `xml:"lemma"`
	Id      string   `xml:"id,attr"`
	Rev     string   `xml:"rev,attr"`
	L       L        `xml:"l"`
	F       []F      `xml:"f"`
}

type L struct {
	T string `xml:"t,attr"`
}

type F struct {
	T string `xml:"t,attr"`
}

type Address struct {
	City  string `xml:"city" json:"city,omitempty"`
	State string `xml:"state" json:"state,omitempty"`
}

func main() {
	/*
		bs2, err := ioutil.ReadFile("test2.txt")
		if err != nil {
			fmt.Println("Error while opening a file #2...")
			return
		}
		str2 = string(bs2)
		println("Content #2:", str2)
	*/

	strToSplit := `Всекитайское собрание народных представителей 11 марта одобрило изменения в конституцию, которые позволяют председателю КНР и его заместителю оставаться у власти неограниченное количество сроков, сообщает BBC.
	В голосовании принимали участие 2964 депутата. Среди них изменения в конституцию поддержали 2959 человек; двое проголосовали против, еще трое воздержались.	
				Отменить пункт конституции, который ограничивает власть главы КНР и его заместителя двумя пятилетними сроками, в конце февраля 2018 года предложила Коммунистическая партия Китая.
	Нынешний председатель КНР Си Цзиньпин возглавил страну в 2013 году. По действующей конституции ему пришлось бы покинуть свой пост не позже 2023 год`

	r := strings.NewReplacer(".", " ", ",", " ", ";", " ")
	strToSplit = r.Replace(strToSplit)
	fmt.Printf("%q \n", strings.Fields(strToSplit))

	//rawXmlData := "<data><person><firstname>Nic</firstname><lastname>Raboy</lastname><address><city>San Francisco</city><state>CA</state></address></person><person><firstname>Maria</firstname><lastname>Raboy</lastname></person></data>"
	rawXmlData, _ := strconv.Unquote(setData())
	//rawXmlData := setData()
	//fmt.Printf(rawXmlData)
	var data Lemmata
	xml.Unmarshal([]byte(rawXmlData), &data)
	fmt.Println(data.LemmList)
	//jsonData, _ := json.Marshal(data)
	//fmt.Println(string(jsonData))
}

func setData() (str string) {

	str = strconv.Quote(`
    <lemmata>
        <lemma id="1" rev="402007">
			<l t="абажур">
				<g v="NOUN"/>
				<g v="inan"/>
				<g v="masc"/>
			</l>
			<f t="абажур">
				<g v="sing"/>
				<g v="nomn"/>
			</f>
			<f t="абажура">
				<g v="sing"/>
				<g v="gent"/>
			</f>
        </lemma>
    </lemmata>
	`)

	return str
}
