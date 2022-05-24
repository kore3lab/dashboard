import * as d3Force			from "d3-force";

/** 
*/
export namespace HierarchyModel {

	export class Hierarchy extends Map<string, Array<Node>>{}

	export class Metadata {
		name:string
		namespace?:string
		ownerReferences?:Array<Node>
	}

	export class Node {
		name:string
		kind:string
		namespace?:string
		depth:number
		ownerReference?:Node 
		children:Array<Node>

		constructor(kind?:string, metadata?:any) {
			if(metadata) {
				this.kind = kind!
				this.name = metadata.name!
				this.namespace = metadata.namespace
				if(metadata.ownerReferences) {
					this.ownerReference = new Node()
					this.ownerReference.kind = metadata.ownerReferences[0].kind
					this.ownerReference.namespace = metadata.namespace
					this.ownerReference.name = metadata.ownerReferences[0].name
				}
			} else {
				if(kind) this.name = kind;	//argument.lenth==1 then "name"
			}
			this.children  = [];
		}
	}
}


/** 
 * Topology 데이터 모델
*/
export namespace TopologyModel {

	export class Topology {
		public nodes:Node[] = [];
		public links:Link[] = [];
	}


	export class Node implements d3Force.SimulationNodeDatum {
		id:string
		name:string
		kind:NodeKind
		group:string;
		labels?:any
		index?: number
		x?: number
		y?: number
		vx?: number
		vy?: number
		fx?: number | null
		fy?: number | null
	}
	export class Link implements d3Force.SimulationNodeDatum {
		source:any
		target:any
		kind:NodeKind
		hidden:boolean=false
		x?: number
		y?: number
		vx?: number
		vy?: number
		fx?: number | null
		fy?: number | null
	}

	export enum NodeKind {
		SERVICE = "service", POD = "pod", NAMESPACE = "namespace", NODE = "node", CLUSTER = "cluster", CONTAINER = "container",
		USER = "user", GROUP = "group", ROLE = "role", CLUSTER_ROLE = "clusterrole", SERVICE_ACCOUNT = "serviceaccount", ROLEBINDING = "rolebinding", CLUSTER_ROLEBINDING = "clusterrolebinding", SECRET ="secret"
	}
}
