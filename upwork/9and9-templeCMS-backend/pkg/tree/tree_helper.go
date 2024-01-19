package tree

import (
	"context"
	"errors"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
)

type TreeDbHelper interface {
	MetaHelper() database.CrudHelper[database.MysqlCondition, dbNodemeta, int64]
	ValHelper() database.CrudHelper[database.MysqlCondition, dbNodeVal, int64]
	GetNodeWithValue(ctx context.Context, groupID int64) ([]CombinedNode, error)
	AddNode(ctx context.Context, node *CombinedNode) (*CombinedNode, error)
	insertValueForMeta(ctx context.Context, groupID int64, metaRoot *NodeMeta) error
}

type treeHelper struct {
	dbHelper TreeDbHelper
}

func NewTreeHelper(dbHelper TreeDbHelper) treeHelper {
	return treeHelper{dbHelper: dbHelper}
}

func (t *treeHelper) GetMeta(ctx context.Context) (*NodeMeta, error) {
	nodeTypes, err := t.dbHelper.MetaHelper().
		Get(ctx, []string{"id", "parent_id", "name", "type", "required"}, database.NewMysqlConditionHelper())
	if err != nil {
		return nil, err
	}

	metaMap := make(map[int64]*NodeMeta)
	var root *NodeMeta

	for _, nodeType := range nodeTypes {
		metaMap[nodeType.ID] = NewNodeMetaFromDB(nodeType)
		if nodeType.ParentID != nil {
			metaMap[*nodeType.ParentID].AddChild(metaMap[nodeType.ID])
		} else {
			root = metaMap[nodeType.ID]
		}
	}

	if root == nil {
		return nil, errors.New("root field not found")
	}

	return root, nil
}

func (t *treeHelper) GetTree(ctx context.Context, groupID int64) (*Node, error) {
	combinedNodes, err := t.dbHelper.GetNodeWithValue(ctx, groupID)
	if err != nil {
		return nil, err
	}

	nodeMap := make(map[int64]*Node)
	var root *Node

	for _, combinedNode := range combinedNodes {
		nodeMap[combinedNode.V.ID] = NewNodeFromDB(combinedNode.T, combinedNode.V)

		if combinedNode.V.ParentID != nil {
			nodeMap[*combinedNode.V.ParentID].AddChild(nodeMap[combinedNode.V.ID])
		} else {
			root = nodeMap[combinedNode.V.ID]
		}
	}

	return root, nil
}

func (t *treeHelper) SetupValueFromMeta(ctx context.Context, groupID int64) error {
	metaRoot, err := t.GetMeta(ctx)
	if err != nil {
		return err
	}

	return t.dbHelper.insertValueForMeta(ctx, groupID, metaRoot)
}

func (t *treeHelper) AddNode(ctx context.Context, groupID int64, parentID int64, node *Node) (*Node, error) {

	combinedNode := &CombinedNode{
		T: dbNodemeta{
			Name:     node.Name,
			Type:     node.Type,
			Required: node.Required,
		},
		V: dbNodeVal{
			GroupID: &groupID,
			Value:   node.Value.(string), // TODO: need to fix this
		},
	}

	parentNode, err := t.dbHelper.ValHelper().Get(ctx, []string{"id", "parent_id", "group_id"},
		database.NewMysqlConditionHelper().Set("id", database.ConditionOperationEqual, parentID),
	)
	if err != nil {
		return nil, err
	}

	if len(parentNode) != 0 { // if parent id is not there then it is root node
		combinedNode.T.ParentID = &parentNode[0].NodemetaID
		combinedNode.V.ParentID = &parentNode[0].ID
	}

	combinedNode, err = t.dbHelper.AddNode(ctx, combinedNode)
	if err != nil {
		return nil, err
	}

	return NewNodeFromDB(combinedNode.T, combinedNode.V), nil
}
