package main

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	//"io"
	"os"
	"sort"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Must specify input .osm.gz file")
		os.Exit(1)
	}

	inputFile, err := os.Open(os.Args[1])
	handleError(err)
	defer inputFile.Close()

	reader, err := gzip.NewReader(inputFile)
	handleError(err)

	//reader = io.TeeReader(reader, os.Stdout)

	decoder := xml.NewDecoder(reader)
	doProcess(decoder)
}

type Tag struct {
	XMLName xml.Name `xml:"tag"`
	Key     string   `xml:"k,attr"`
	Value   string   `xml:"v,attr"`
}

type TagList []Tag

func (tl TagList) getTag(tagName string) (value string, ok bool) {
	for _, tag := range tl {
		if tag.Key == tagName {
			return tag.Value, true
		}
	}
	return "", false
}

type Node struct {
	XMLName xml.Name `xml:"node"`
	Id      int32    `xml:"id,attr"`
	Tags    TagList  `xml:"tag"`
	Lat     float64  `xml:"lat,attr"`
	Lon     float64  `xml:"lon,attr"`
}
type NodeList []*Node

type MemberXML struct {
	XMLName   xml.Name `xml:"member"`
	Type      string   `xml:"type,attr"`
	Reference int32    `xml:"ref,attr"`
	Role      string   `xml:"role,attr"`
}

type RelationXML struct {
	XMLName xml.Name    `xml:"relation"`
	Id      int32       `xml:"id,attr"`
	Members []MemberXML `xml:"member"`
	Tags    TagList     `xml:"tag"`
}

// A parsed MemberXML structure
const (
	RELATION_TYPE_NODE = iota
	RELATION_TYPE_RELATION
	RELATION_TYPE_WAY
)

type Member struct {
	Type      int
	Reference interface{}
	Role      string
}

type Relation struct {
	Id      int32
	Members []Member
	Tags    TagList
}

// Implement the sort interface
func (nl NodeList) Len() int           { return len(nl) }
func (nl NodeList) Less(i, j int) bool { return nl[i].Id < nl[j].Id }
func (nl NodeList) Swap(i, j int)      { nl[i], nl[j] = nl[j], nl[i] }
func (nl NodeList) findNode(nodeId int32) *Node {
	i := sort.Search(len(nl), func(i int) bool { return nl[i].Id >= nodeId })
	if i < len(nl) && nl[i].Id == nodeId {
		return nl[i]
	} else {
		return nil
	}
}

type OSMData struct {
	Nodes NodeList
	//Relations RelationList

	sortedNodes bool
}

func NewOSMData() (osm *OSMData) {
	osm = new(OSMData)
	osm.Nodes = make(NodeList, 0)
	//osm.Relations = make(RelationList, 0)
	return
}

func (osm *OSMData) findReference(referenceID int32) interface{} {
	if node := osm.Nodes.findNode(referenceID); node != nil {
		return node
	}
	return nil
}

func (osm *OSMData) sortNodes() {
	if !osm.sortedNodes {
		sort.Sort(osm.Nodes)
	}
	osm.sortedNodes = true
}

func doProcess(decoder *xml.Decoder) {
	var (
		err   error
		token xml.Token
		osm   = NewOSMData()
	)
	for {
		if token, err = decoder.Token(); err != nil {
			break
		}

		if se, ok := token.(xml.StartElement); ok {
			switch se.Name.Local {
			case "node":
				if node := decodeNode(decoder, &se); node != nil {
					osm.Nodes = append(osm.Nodes, node)
				}
			case "way":
				osm.sortNodes() // If we get to this point, it means we've finished getting the list of nodes
				continue
			case "relation":
				osm.sortNodes() // If we get to this point, it means we've finished getting the list of nodes
				decodeRelation(osm, decoder, &se)
			}
		}
	}
}

func decodeNode(decoder *xml.Decoder, startElement *xml.StartElement) (node *Node) {
	if err := decoder.DecodeElement(&node, startElement); err != nil {
		return nil
	}

	if len(node.Tags) == 0 {
		return nil
	}
	return node
}

func decodeRelation(osm *OSMData, decoder *xml.Decoder, startElement *xml.StartElement) (relation *Relation) {
	var relationXML RelationXML
	if err := decoder.DecodeElement(&relationXML, startElement); err != nil {
		return
	}

	// Convert the relationXML to a proper Relation object
	relation = new(Relation)
	relation.Id = relationXML.Id
	relation.Tags = relationXML.Tags
	relation.Members = make([]Member, 0, len(relationXML.Members))

	for _, memberXML := range relationXML.Members {
		//fmt.Println(memberXML.Role)
		if ref := osm.findReference(memberXML.Reference); ref != nil {
			var relationType int = -1
			switch memberXML.Type {
			case "node":
				relationType = RELATION_TYPE_NODE
			case "relation":
				relationType = RELATION_TYPE_RELATION
			case "way":
				relationType = RELATION_TYPE_WAY
			default:
				panic("Unsupported relation type found: " + memberXML.Type)
			}

			member := Member{
				Type:      relationType,
				Reference: ref,
				Role:      memberXML.Role,
			}
			relation.Members = append(relation.Members, member)
		}
	}

	return
}
