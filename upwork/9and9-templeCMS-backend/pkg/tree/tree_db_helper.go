package tree

import (
	"context"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
)

type CombinedNode struct {
	T dbNodemeta
	V dbNodeVal
}

type treeDbHelper struct {
	metaHelper *dbNodemetaHelper
	valHelper  *database.BaseHelper[dbNodeVal]
	db         database.Connection

	configID int64
}

func NewTreeDbHelper(db database.Connection, valTable string, configID int64) TreeDbHelper {
	return &treeDbHelper{
		metaHelper: NewDbNodemetaHelper(db, configID),
		valHelper:  NewDbNodeValHelper(db, valTable),
		db:         db,

		configID: configID,
	}
}

func (h *treeDbHelper) MetaHelper() database.CrudHelper[database.MysqlCondition, dbNodemeta, int64] {
	return h.metaHelper
}

func (h *treeDbHelper) ValHelper() database.CrudHelper[database.MysqlCondition, dbNodeVal, int64] {
	return h.valHelper
}

func (h *treeDbHelper) AddNode(ctx context.Context, node *CombinedNode) (*CombinedNode, error) {
	tx, err := database.NewTransationFromConnection(ctx, h.db)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	metaHelper, err := h.metaHelper.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	node.T.ConfigID = h.configID

	nodeMeta, err := metaHelper.Create(ctx, &node.T)
	if err != nil {
		return nil, err
	}

	node.V.NodemetaID = node.T.ID
	valHelper, err := h.valHelper.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	nodeVal, err := valHelper.Create(ctx, &node.V)
	if err != nil {
		return nil, err
	}

	return &CombinedNode{T: *nodeMeta, V: *nodeVal}, nil
}

func (t *treeDbHelper) insertValueForMeta(ctx context.Context, groupID int64, metaRoot *NodeMeta) error {
	tx, err := database.NewTransationFromConnection(ctx, t.db)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	valHelper, err := t.valHelper.NewWithTx(tx)
	if err != nil {
		return err
	}

	root, err := valHelper.Create(ctx, &dbNodeVal{
		NodemetaID: metaRoot.ID,
		Value:      "",
		GroupID:    &groupID,
	})
	if err != nil {
		return err
	}

	err = t.addTreeNodes(ctx, groupID, root.ID, metaRoot.Children, valHelper)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (t *treeDbHelper) addTreeNodes(ctx context.Context, groupID int64, parentID int64, nodes []*NodeMeta, valHelper *database.BaseHelper[dbNodeVal]) error {
	for _, meta := range nodes {
		resp, err := valHelper.Create(ctx, &dbNodeVal{
			ParentID:   &parentID,
			NodemetaID: meta.ID,
			Value:      "",
			GroupID:    &groupID,
		})
		if err != nil {
			return err
		}

		if meta.Type == Array {
			continue
		}

		err = t.addTreeNodes(ctx, groupID, resp.ID, meta.Children, valHelper)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *treeDbHelper) GetNodeWithValue(ctx context.Context, groupID int64) ([]CombinedNode, error) {

	tableName := h.valHelper.GetTableName()
	query := `SELECT tree_fields.id as field_id, 
			tree_fields.parent_id as field_parent_id,
			name as field_name,
			type as field_type,
			required as field_required,

			` + tableName + `.parent_id as value_parent_id, 
			` + tableName + `.id as value_id, 
			` + tableName + `.value as value
		FROM tree_fields
		LEFT JOIN ` + tableName + ` on tree_fields.id = ` + tableName + `.node_type_id
		WHERE tree_fields.config_id = ? and ` + tableName + `.group_id = ?`

	return database.QueryScanner(ctx, h.db, func() func(*CombinedNode) []interface{} {
		return func(c *CombinedNode) []interface{} {
			return []interface{}{
				&c.T.ID,
				&c.T.ParentID,
				&c.T.Name,
				&c.T.Type,
				&c.T.Required,
				&c.V.ParentID,
				&c.V.ID,
				&c.V.Value,
			}
		}
	}(), query, h.configID, groupID)

}
