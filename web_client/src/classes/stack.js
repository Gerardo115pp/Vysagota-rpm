class StackNode {
  constructor(value) {
    this.value = value;
    this.next = null;
  }
}

class Stack {
    constructor(limit=Infinity) {
        this.limit = limit;
        this.top = null;
        this.length = 0;
    }

    push = (value) => {
        if (this.length < this.limit) {
            const node = new StackNode(value);
            node.next = this.top;
            this.top = node;
            this.length++;
            return node;
        }
        return null;
    }
    
    pop = () => {
        let node = this.top;
        if (node) {
            this.top = node.next;
            this.length--;
            return node.value;
        } else {
            return null; 
        }
    }

    peek = () => this.top != null ? this.top.value : null;

    isEmpty = () => this.length === 0;
}

export default Stack;