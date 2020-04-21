import * as d3Array	from "d3-array";
import {Kiali}		from "../../model/graph.model";
import * as model	from "./model";

/**
 * Kiali "app" 그래프 데이터 정제 처리 클래스
 * 
 */
export class KialiParser {

	/**
	 * Kiali 원본 데이터를 그래프 그리기 위한 데이터로 변환하는 작업
	 * 		model.Kiali.Graph -> model.Source
	 */
	public static parse(graph:Kiali.Graph):model.Source {

		let nodeMap:Map<string,model.Node> = new Map<string,model.Node>();
		let edgeList:Array<model.Edge> = [];
		let nodeList:Array<model.Node> = [];

		
		// 노드 데이터
		if(!graph.elements || !graph.elements["nodes"] || !graph.elements["edges"]) return;
		nodeList = graph.elements.nodes.map( (d:{data:Kiali.Node}) =>{
			d.data.id =`_${d.data.id}`;
			if(d.data.parent)  d.data.parent =`_${d.data.parent}`;
			d.data.isRoot = (d.data.isRoot || d.data.isUnused);
			return d.data;
		});
		// 노드 데이터  : id 를 키로  노드정보를 key-value 형태로 변경
		nodeMap = nodeList.reduce((accumulator:Map<string,model.Node>, d:model.Node) => {
				accumulator.set(d.id, d);
				return accumulator;
			}, new Map());

		// 엣지 데이터		
		graph.elements.edges.forEach( (d:any) => {
			if(d.data.source==d.data.target) return;

			d.data.id =`_${d.data.id}`;
			d.data.source =`_${d.data.source}`;
			d.data.target =`_${d.data.target}`;
			// 양방향 
			const idx:number = edgeList.findIndex((dd:model.Edge)=>(dd.target==d.data.source && dd.source==d.data.target));
			if(idx>-1) {
				edgeList.splice(idx, 1);
				d.data.isTwoWay = true;
			}
			edgeList.push(d.data);
		});

		// // 엣지 데이터  : source 를 기준으로 edges 정보를  map 형태로 변경
		let sourceEdge:model.EdgeSource = {
			data: edgeList,
			map: d3Array.group<model.Edge,string>(edgeList, d=>d.id),
			source: d3Array.group<model.Edge,string>(edgeList, d=>d.source),
			target: d3Array.group<model.Edge,string>(edgeList, d=>d.target)
		}

		// parent속성으로 그룹핑 
		let groupMap:Map<string,Array<model.Node>> = d3Array.group<model.Node,string>(nodeList, d=>d.parent);

		// 그룹 리스트를 만든다.
		// 		- undefined 이면 isGroup 이거나 parent 없는 단독 group 임
		let groupList:Map<string, model.Group> = new Map<string, model.Group>();
		if(groupMap.get(undefined)) {
			groupMap.get(undefined)
				.forEach((d:model.Node)=>{

					let group:model.Group = {
						id: d.id,
						name: d.app, 
						namespace: d.namespace,
						app: d.app,
						groups:[],
						nodes: d.isGroup ? groupMap.get(d.id): [d]
					};

					groupList.set(group.id, group);	//그룹 리스트

				});
		}

		// 재귀호출하여 parent-child 그룹으로 재정렬 작업
		groupList.forEach((g:model.Group)=>{
			this.addGroup(g, groupList, nodeMap, sourceEdge.source);
		});

		// 정렬된 결과 list 로 담는다.
		let source:model.Source = new model.Source();
		groupList.forEach((g:model.Group, key:string)=>{
			source.groups.push(g);
		});

		source.edges = sourceEdge;
		return source;
	}

	/**
	 * parent-child 그룹으로 재정렬 작업
	 * 
	 */
	private static addGroup(group:model.Group, groupList:Map<string, model.Group>, nodesMap: Map<string,model.Node>, edgeMap:Map<string,model.Edge[]>) {

		const grpIds:Array<string> = this.getChildrenGroupIds(group, nodesMap, edgeMap);

		grpIds.forEach((d:string)=>{
			const g:model.Group = groupList.get(d);
			if(g) {
				this.addGroup(g, groupList, nodesMap, edgeMap);
				group.groups.push(g);	//그룹추가
				groupList.delete(d);	//그룹리스트에서 제외
			}
		});

	}

	/**
	 * edge 정보를 사용하여 Child 그룹 id 얻기
	 * 
	 */
	private static getChildrenGroupIds(g:model.Group, nodesMap: Map<string,model.Node>, edgeMap:Map<string,model.Edge[]>):Array<string> {

		let ids:Array<string> = [];

		g.nodes.forEach((d:model.Node)=>{
			const edges:Array<model.Edge> = edgeMap.get(d.id);
			if(edges) {
				edges.forEach((e:model.Edge)=> {
					const id:string = this.getGroupKey(e.target, nodesMap)
					if(id && g.id != id && !ids.includes(id) ) {
						ids.push(id)
					}
				});
			}
		});

		return ids;

	}



	/**
	 * 해당 노드의 그룹키 리턴
	 * 그룹키는 app.namespace 또는 service.namespace
	 */
	private static getGroupKey(id, nodesMap: Map<string,model.Node>):string {
		const node:model.Node = nodesMap.get(id);
		return (node.parent) ? node.parent: node.id;
	}

}
