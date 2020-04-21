import * as d3Scale		from "d3-scale"
import {ConfigModel}	from "../../model/models";
import {Transform}		from "../../utils/transform";
import * as model		from "./model";

export class GraphUI {


	/**
	 * 해당 element 와 관계있는 노드&선 리턴
	 * 		- edges 정보를 활용
	 * 
	 * @param element	시작 element
	 * @param mapEdges 연결정보
	 */
	public static getRelationNodes(element:SVGElement, sourceEdge:model.EdgeSource):Map<string, SVGElement> {

		let id:string = element.getAttribute("id")?element.getAttribute("id"):"this";
		let ids:Array<string> = [id];

		element.querySelectorAll("g.node").forEach((e:SVGElement) => ids.push(e.getAttribute("id")));


		let map:Map<string, SVGElement> = GraphUI.getRelationNodesById(element instanceof SVGSVGElement?<SVGSVGElement>element: element.ownerSVGElement, ids, sourceEdge);
		if(!map.has(id)) map.set(id, element);

		return map;


	}

	/**
	 * 해당 아이디의 element 와 관계있는 노드&선 리턴
	 * 		- edges 정보를 활용
	 * 
	 * @param svg	SVG element
	 * @param id	element 아이디
	 * @param mapEdges 연결정보
	 */
	public static getRelationNodesById(svg:SVGSVGElement, ids:Array<string>, sourceEdge:model.EdgeSource):Map<string, SVGElement> {

		let map:Map<string, SVGElement> = new Map();

		ids.forEach( (id:string) => {

			// source 기준 edge map
			if(sourceEdge.source) {
				let edges:model.Edge[] = sourceEdge.source.get(id);
				edges && edges.forEach( (value:model.Edge) => {
					if(value.id && !map.has(value.id)) map.set(value.id, svg.querySelector(`#${value.id}`));				// line (out to)
					if(value.target && !map.has(value.target)) map.set(value.target, svg.querySelector(`#${value.target}`));	//target node
				});
			}

			// target 기준 edge map
			if(sourceEdge.target) {
				let edges = sourceEdge.target.get(id);
				edges && edges.forEach( (value:model.Edge) => {
					if(value.id && !map.has(value.id)) map.set(value.id, svg.querySelector(`#${value.id}`));			//line (into)
					if(value.source && !map.has(value.source)) map.set(value.source, svg.querySelector(`#${value.source}`));	//source node
				});
			}

			// 라인인 경우 edge map
			if(sourceEdge.map) {
				let value:model.Edge= sourceEdge.map[id];
				if(value) {
					if(value.source && !map.has(value.source)) map.set(value.source, svg.querySelector(`#${value.source}`));	//source node
					if(value.target && !map.has(value.target)) map.set(value.target, svg.querySelector(`#${value.target}`));	//target node
				}
			}

		});

		return map;

	}

	/**
	 * Edges x1, x2, y1,y2 위치 정의
	 */
	public static calcEdgesPosition(outlineEl:d3.Selection<SVGGElement,any,SVGElement,any>, sourceEdge:model.EdgeSource, options:ConfigModel.Mesh, bounds:DOMRect) {
		
		
		const edges:Array<model.Edge> = sourceEdge.data;

		// 위치보정 source 가 target 좌측에 있다면 위치 조정
		// let out:boolean = false
		// while(!out) {
		// 	out = edges.every((d:model.Edge)=>{
		// 		try {
		// 			let source:SVGElement =  outlineEl.select<SVGElement>(`#${d.source} .ico`).node();
		// 			let target:SVGElement =  outlineEl.select<SVGElement>(`#${d.target} .ico`).node();
					
		// 			if(!source || !target) return true;
		// 			if(d.isTwoWay) return true;	//양뱡향인 경우는 제외

		// 			if(source && target) {
		// 				let sourceRect:ClientRect =  source.getBoundingClientRect();
		// 				let targetRect:ClientRect =  target.getBoundingClientRect();

		// 				let x1 = sourceRect.right - bounds.left,
		// 				x2 = targetRect.left - bounds.left;

		// 				if(x1 >= x2) {
		// 					// target이 좌측에 있다면 target source 의 오른쪽으로 이동
		// 					let el:HTMLElement = target.parentElement.parentElement.parentElement;
		// 					if(el) {
		// 						Transform.instance(el).shiftX(x1+options.node.distance.x);
		// 					}
		// 					return false;
		// 				}
		// 			}
		// 			return true;

		// 		} catch (e) {
		// 			return false;
		// 		}
		// 	})
		// }

		// 선 위치 x1,x2,y1,y2 값 정의

		edges.forEach( (d:model.Edge) => {
			let source:SVGElement =  outlineEl.select<SVGElement>(`#${d.source} .ico`).node();
			let target:SVGElement =  outlineEl.select<SVGElement>(`#${d.target} .ico`).node();
			
			if(!source || !target) return;

			// 양뱡이고 source 가 target 우측에 있으면 source <-> target
			if(d.isTwoWay && source.getBoundingClientRect().left > target.getBoundingClientRect().left ) {
				source =  outlineEl.select<SVGElement>(`#${d.target} .ico`).node();
				target =  outlineEl.select<SVGElement>(`#${d.source} .ico`).node();
			}

			const sourceRect:ClientRect =  source.getBoundingClientRect();
			const targetRect:ClientRect =  target.getBoundingClientRect();

			// 시작x 도착x 보정처리위한 값
			const offset = {
				source: source.classList.contains("service")? -options.node.r/2: 0,
				target: target.classList.contains("service")? 0: -4
			}

			d.x1 = sourceRect.right - bounds.left + offset.source,
			d.x2 = targetRect.left - bounds.left + offset.target,
			d.y1 = sourceRect.top + sourceRect.height/2 - bounds.top,
			d.y2 = targetRect.top + targetRect.height/2 - bounds.top

		});

	}


}


export class Traffic {

	//edge traffic 에서 response 에러율 구하기
	public static getErrorRatio(traffic:model.Traffic):number {

		const keys:string[]  = Object.keys(traffic.responses);
		let ratio:number = 0;
		
		keys.forEach(key=> {
			if(traffic.responses[key].flags["FI"]) ratio += traffic.responses[key].flags.FI;
		});
		return ratio;
	}

	public static getStatus(d:model.Edge, conf:ConfigModel.Config):""|"error"|"error anmi"|"warn"|"success"|"idle" {

		if(d.traffic.responses) {
			let errorRatio:number = Traffic.getErrorRatio(d.traffic);
			if(errorRatio/100 > conf.global.health.error.ratio) return conf.mesh.alert.error.animation.edge ? "error anmi": "error";
			else if(errorRatio/100 > conf.global.health.warn.ratio) return "warn";
			else return "success";
		} else return "idle";

	}

}
