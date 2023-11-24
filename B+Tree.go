package main

import "sync"

// B树上所有节点中孩子节点个数的最大值称为阶
// 每个节点最多有m个子树，每个节点上的关键字最多有m-1个关键字
// 除根节点外，其他的每个分支至少有ceil(M/2)个子树，至少含有ceil(M/2)-1个关键字。
// 每个节点中的关键字都按照大小顺序排列，每个关键字的左子树的所有关键字都小于它，每个关键字的右子树都大于它。
// 所有的叶子节点都位于同一层

// B+树的一些特性
// 每个节点这至少有M个子树
// 除根节点之外，每个节点至少有ceil(m/2)个子树
// 结点的子树等于关键字个数
// 叶子节点包含了所有关键字的信息，以及指向含这些关键记录的指针，并且叶子节点本身按照关键字的大小来顺序链接
// 所有的非叶子节点，看成索引部分，结点中仅含有其子树中的最大关键字

// B+树和B树相对，主要有以下区别
//非叶子节点只存储键值信息，数据记录都存放在叶子节点中。
//所有叶子节点之间都有一个链指针。
//非叶子节点的关键字的个数与其子树的个数相同，不像B树，子树的个数总比关键字个数多1个。

//B+树通常用于数据库索引，例如Mysql的InnoDB存储引擎以及MyISAM存储引擎的索引文件中使用的就是B+树。
//一般来说，数据库的索引都比较大，不可能全部存储在内存中，因此索引往往以文件的形式存储的磁盘上。
//系统从磁盘读取文件到内存时，需要定位到文件所在的位置：文件所在柱面号，磁盘号，扇区号。这个操作时非常耗时的，远高于内存操作。
//考虑到磁盘IO是非常高昂的操作，操作系统在读取文件时做了一些优化，系统从磁盘读取文件到内存时是以磁盘块（block）为基本单位的，
//位于同一个磁盘块中的数据会被一次性读取出来，而不是需要什么取什么。
//每一次IO读取的数据我们称之为一页(page)。具体一页有多大数据跟操作系统有关，一般为4k或8k。

//由于磁盘IO非常耗时，因此评价一个数据结构作为索引的优劣最重要的指标就是在查找过程中是否能够有效减少磁盘I/O的操作次数。
//Mysql选择使用B+树作为索引文件的数据结构，主要基于B+树的以下特点：

//B+树的磁盘读写代价更低 B+树的内部结点只有关键字，没有数据，一个结点可以容纳更多的关键字。
//查询时一次性读入内存中的关键字也就越多，相对来说I/O读写次数也就降低了。
//B+树查询效率更加稳定 B+树内部结点不保存数据，而只是叶子结点中数据的索引。
//所以任何关键字的查找必须走一条从根结点到叶子结点的路。所有关键字查询的路径长度相同，导致每一个数据的查询效率相当。
//B+树便于范围查询 所有叶子节点形成有序链表，对于数据库中频繁使用的范围查询，B+树有着更高的性能。。
//在InnoDB中，表数据文件本身就是按B+树组织的一个索引结构，它使用数据库主键作为Key，叶节点保存了完整的数据记录。
//InnoDB中有页（Page）的概念，页是其磁盘管理的最小单位。InnoDB中默认每个页的大小为16KB，
//可通过参数innodb_page_size将页的大小设置为4K、8K、16K。
//InnoDB中，B+Tree的高度一般都在2~4层。由于根节点常驻内存的，因此查找某一键值的行记录时最多只需要1~3次磁盘I/O操作。
//因为InnoDB的数据文件本身要按主键聚集，所以InnoDB要求表必须有主键（MyISAM可以没有），
//如果没有显式指定，则MySQL系统会自动选择一个可以唯一标识数据记录的列作为主键，
//如果不存在这种列，则MySQL自动为InnoDB表生成一个隐含字段作为主键，这个字段长度为6个字节，类型为长整形。
//聚集索引这种实现方式使得按主键的搜索十分高效，
//但是辅助索引搜索需要检索两遍索引：首先检索辅助索引获得主键，然后用主键到主索引中检索获得记录。

// B+树中插入数据
// 1) 首先要定位数据所在的叶子节点，然后将数据插入到该节点，出入数据后不能破坏关键字的排列顺序
// 2) 如果插入的元素之后关键字数目小于等于阶数M，那么直接完成插入操作
// 3) 如果插入的元素为该节点的最大值，那么需要修改其父节点中索引值
// 4) 若插入元素后该节点关键字数目>阶数M，则需要将该结点分裂为两个结点，关键字的个数分别为：floor((M+1)/2)和ceil((M+1)/2)。
// 5) 若分裂结点后导致父节点的关键字数目>阶数M，则父节点也要进行相应的分裂操作。

// B+树中删除数据
// 当删除某结点中最大或者最小的关键字，就会涉及到更改其双亲节点一直到根节点中所有索引值的更改
// 删除关键字后，如果当前节点中关键字个数大于M/2 直接删除
// 在删除关键字后，如果导致其结点中关键字个数<[M/2]，若其兄弟结点中含有多余的关键字，可以从兄弟结点中借关键字。
// 在删除关键字后，如果导致其结点中关键字个数<[M/2]，并且其兄弟结点没有多余的关键字，则需要同其兄弟结点进行合并。
// 结点合并后，需要修改父结点的关键字的个数，若父结点的关键字个数<[M/2]，需要依照以上规律进行处理。

//BPItem用于数据记录。
//MaxKey：用于存储子树的最大关键字
//Nodes：结点的子树（叶子结点的Nodes=nil）
//Items：叶子结点的数据记录（索引结点的Items=nil）
//Next：叶子结点中指向下一个叶子结点，用于实现叶子结点链表

// 叶子存储的内容
type BPItem struct {
	Key int64
	Val interface{}
}

type BPNode struct {
	MaxKey int64     // 当前节点的最大值
	Nodes  []*BPNode // 节点的子树（叶子节点为nil）
	Items  []BPItem  // 叶子节点存储的内容，索引节点为nil
	Next   *BPNode   // 叶子节点中使用的链表
}

// 在叶子节点上插入一个值 如果遇到索引相同的值，直接覆盖
func (node *BPNode) setValue(key int64, value interface{}) {
	item := BPItem{key, value}
	num := len(node.Items)
	if num < 1 {
		node.Items = append(node.Items, item)
		node.MaxKey = key
		return
	} else if key < node.Items[0].Key {
		node.Items = append([]BPItem{item}, node.Items...) // 在前面插入一个元素
		return
	} else if key > node.Items[num-1].Key {
		node.Items = append(node.Items, item)
		node.MaxKey = key
		return
	}

	// 找到需要插入的点
	for i := 0; i < num; i++ {
		if node.Items[i].Key > key {
			node.Items = append(node.Items, BPItem{}) // 为了给
			copy(node.Items[i+1:], node.Items[i:])
			node.Items[i] = item
			return
		} else if node.Items[i].Key == key {
			node.Items[i] = item
			return
		}
	}

}

type Btree struct {
	mutex sync.RWMutex
	root  *BPNode
	M     int // 记录最大阶数
	halfM int // 记录ceil(M/2)
}

func NewBtree(M int) *Btree {
	if M < 3 {
		M = 3
	}

	btree := &Btree{}
	btree.M = M
	btree.halfM = (M + 1) / 2
	btree.root = NewLeafNode(M)

	return btree
}

// 申请width+1 是因为插入可能暂时出现节点key大于width的情况，等待后期再分裂处理
func NewLeafNode(width int) *BPNode {
	node := &BPNode{}
	node.Items = make([]BPItem, width+1)
	node.Items = node.Items[0:0]
	return node
}

func (t *Btree) Get(key int64) interface{} {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	node := t.root
	// 从btree上
	for i := 0; i < len(node.Nodes); i++ {

		if key <= node.Nodes[i].MaxKey {
			node = node.Nodes[i]
			i = 0
		}
	}

	// 如果没有达到叶子节点
	if len(node.Nodes) > 0 {
		return nil
	}

	// 进入叶子节点中查找
	for i := 0; i < len(node.Items); i++ {
		if node.Items[i].Key == key {
			return node.Items[i].Val
		}
	}
	return nil
}

func (t *Btree) Set(key int64, value interface{}) {
	t.mutex.Lock()
	t.mutex.Unlock()

	// 首先加入之后再判断

}

//func (t *Btree) splitNode(node *BPNode) *BPNode {
//	if len(node.Nodes) > t.M {
//		// 创建新节点
//		halfw := t.M/2 + 1
//		node2 := NewLeafNode(t.M)
//		node2.Nodes=append(node2.Nodes,node)
//
//	} else if len(node.Items) > t.M {
//
//	}
//
//}

func setValue(parent *BPNode, node *BPNode, key int64, value interface{}) {

	//找到插入的叶子节点位置？
	for i := 0; i < len(node.Nodes); i++ {
		if key <= node.Nodes[i].MaxKey {

		}
	}
}
