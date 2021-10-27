package main

import (
	"fmt"
	"github.com/sjk00000001/gotype"
)

func Reverse (node *gotype.LNode){
	if node==nil || node.Next ==nil {return}
	var pre *gotype.LNode		//new()或&实例化后得到的指针初值不是nil,var声明的指针初值为nil；interface{}类型初值也是nil
	cur := node.Next
	var next *gotype.LNode
	for cur != nil {
		next =cur.Next
		cur.Next = pre			//第一个结点的后继结点需要赋值nil，所以第一次循环一定要让当前结点为第一个结点，利用上pre最开始的初值nil
		pre = cur				//否则就要想办法设立if条件来为第一个结点的后继结点赋值nil
		cur = next
	}
	node.Next = pre   			//头结点的后继结点可以直接赋值最后一次循环的cur也就是pre的值
}

func main (){
	head :=&gotype.LNode{} //head肯定不是nil，保证了该链表一定是带头结点的单链表
	fmt.Println("就地逆序")
	gotype.CreateNode(head, 8)
	gotype.PrintNode("逆序前：", head)
	Reverse(head)
	gotype.PrintNode("逆序后：", head)
	fmt.Scanln()
}
