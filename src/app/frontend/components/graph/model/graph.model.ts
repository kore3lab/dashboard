import * as d3Force		from "d3-force";

/** 
 * Kiali 데이터 모델
*/
export namespace Kiali {

	export class Graph {
		timestamp:number
		duration:number
		graphType:string
		elements: {
			nodes: Array<{data:Node}>
			edges: Array<{data:Edge}>
		}
	}

	export class Node {
		id:string
		namespace:string
		nodeType:"app"|"service"|"workload"
		app:string			// app 명
		service?:string		// service 명, nodeType="service" 인 경우만 존재
		isGroup:string		// 그룹 , "app"
		isRoot?:boolean		// ㅅ
		isOutside?:boolean	// external
		isInaccessible:boolean	//restrict node
		isUnused?:boolean	// 미사용 노드 여부 (unusedNode 조회시)
		workload?:string	// 워크로드명 node.app + node.version (예:productpage-v1)
		version?:string		// 버전
		parent?:string		// parent 아이디
		hasVS?:boolean		// VirtualService 여부
		hasMissingSC:boolean// missing sidecar 여부
		hasCB:boolean		// circuit breaker 여부
		health?:Health		// Custom 필드
		isServiceEntry?:string	//service entry	(MESH_EXTERNAL)
		traffic:Array<Kiali.Traffic>
	}
	export class Edge {
		id:string
		source:string
		target:string
		traffic: Traffic
		isMTLS?:number			//mTLS 성공%
		responseTime?:number;	//response time
	}

	export class Traffic {
		protocol:string
		rates: {
			http?:number
			http5xx?:number
			httpPercentErr?:number
			httpPercentReq?:number	// 요청 퍼센트 (weighted routing 일 경우)
			httpIn?:number
			httpIn5xx?:number
			httpOut?:number
			tcpOut?:number
		}
		responses?: {
			[statusCode:string] : {
				flags: {
					"-"?:number
					FI?:number
				}
				hosts: {}
			}
		}
	}

	// Health
	export class HealthDictionary {
		[key:string] : Health
	}
	export class Health {
		requests:Requests							// request 정보
		workloadStatuses?:Array<WorkloadStatuse>	// subset정보 (nodeType = "app" 인 경우만)
	}


	// Health - nodeType = ("app", "service","workload") request 정보
	export class Requests {
		errorRatio:number
		inboundErrorRatio:number
		outboundErrorRatio:number
	}

	// Health - nodeType = "app" subset정보
	export class WorkloadStatuse {
		name:string
		desiredReplicas:number
		currentReplicas:number
		availableReplicas:number
	}

}
/** 
 * Topology 데이터 모델
*/
export namespace Tree {

	export class Tree {
		name:string
		kind:NodeKind
		children?:Array<Tree>
	}

	export enum NodeKind {
		USER = "user", GROUP = "group", ROLE = "role", CLUSTER_ROLE = "clusterrole", SERVICE_ACCOUNT = "serviceaccount", ROLEBINDING = "rolebinding", CLUSTER_ROLEBINDING = "clusterrolebinding", SECRET ="secret"
	}

}

/** 
 * Topology 데이터 모델
*/
export namespace Topology {

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
		SERVICE = "service", POD = "pod", NAMESPACE = "namespace", NODE = "node", CLUSTER = "cluster",
		USER = "user", GROUP = "group", ROLE = "role", CLUSTER_ROLE = "clusterrole", SERVICE_ACCOUNT = "serviceaccount", ROLEBINDING = "rolebinding", CLUSTER_ROLEBINDING = "clusterrolebinding", SECRET ="secret"
	}
}
