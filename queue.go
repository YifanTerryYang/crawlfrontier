package crawlfrontier

type Queue struct {
	head *QueueNode
	tail *QueueNode
	length int
}

type QueueNode struct {
	next *QueueNode
	body *URLEntity
}

type URLEntity struct {
	url string
	searchcriteria string
}

func (q *Queue) Add(entity *URLEntity) {
	newnode := &QueueNode{body:entity}
	if q.head == nil {
		q.head = newnode
		q.tail = newnode
		q.length = 1
	} else {
		q.tail.next = newnode
		q.tail = newnode
		q.length += 1
	}
}

func (q *Queue) Get() (result *URLEntity) {
	if q.head == nil {
		result = nil
	} else if q.head == q.tail {
		result = q.head.body
		q.head = nil
		q.tail = nil
		q.length = 0
	} else {
		result = q.head.body
		q.head = q.head.next
		q.length -= 1
	}

	return
}

func (q *Queue) Size() int {
	return q.length
}