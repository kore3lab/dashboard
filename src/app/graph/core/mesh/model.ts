import {Kiali}	from "../../model/graph.model";

/** 
 * 서비스 메시 그래프 데이터 모델
*/
export class Group  {
	id:string
	name:string
	namespace?:string
	app?:string
	nodes:Array<any> = []
	groups:Array<Group> = []
	// isRoot?:boolean				// ROOT 
	isOutline?:boolean			// 테두리 display 여부
}
export class Node extends Kiali.Node {
	// name?:string
	// requests?: Kiali.Requests 
	// workloadStatuses?: Array<Kiali.WorkloadStatuse> = []
}
export class Edge extends Kiali.Edge {
	x1?:number;
	x2?:number;
	y1?:number;
	y2?:number;
	isTwoWay?:boolean;
}

export class Source {
	constructor() {
		this.groups = [], this.edges = new EdgeSource()
	}
	groups:Array<Group>
	edges:EdgeSource
}

export class EdgeSource {
	data:Array<Edge> = [];
	map:Map<string,Edge[]> = new Map();
	source:Map<string,Edge[]> = new Map();
	target:Map<string,Edge[]> = new Map();
}

export class Traffic extends Kiali.Traffic{
}
