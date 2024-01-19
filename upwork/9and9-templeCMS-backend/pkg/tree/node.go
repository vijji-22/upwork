package tree

type NodeMeta struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Type     ValueType `json:"type"`
	Required bool      `json:"required"`

	Children []*NodeMeta `json:"children,omitempty"`
}

func NewNodeMetaFromDB(nodeType dbNodemeta) *NodeMeta {
	return &NodeMeta{
		ID:       nodeType.ID,
		Name:     nodeType.Name,
		Type:     nodeType.Type,
		Required: nodeType.Required,
	}
}

func (n *NodeMeta) AddChild(child *NodeMeta) {
	n.Children = append(n.Children, child)
}

type Node struct {
	ID int64 `json:"id"`
	NodeMeta
	Value interface{} `json:"value"`

	Children []*Node `json:"children,omitempty"`
}

func NewNodeFromDB(nodeType dbNodemeta, nodeVal dbNodeVal) *Node {
	return &Node{
		ID:       nodeVal.ID,
		NodeMeta: *NewNodeMetaFromDB(nodeType),
		Value:    nodeVal.Value,
	}
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}
