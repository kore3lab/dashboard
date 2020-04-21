import * as d3			from "d3";
import * as d3Select	from "d3-selection";
import {ConfigModel}	from "../../model/models";
import * as model		from "./model";
import {GraphUI}		from "./lang";

/**
 * 이벤트 처리 클래스
 * 
 */
export class EventBinder {

	/**
	 * 마우스 over / out 시 하일라이트 처리
	 * 
	 * @param el 대상 element
	 * @param options  옵션
	 */
	public static onHighlight(element:d3.Selection<SVGElement,any,SVGElement,any>, sourceEdge:model.EdgeSource, options:ConfigModel.Mesh) {

		if(!element || !element.node()) return;

		let isNodes:boolean = (element.node().tagName == "rect");	//rect nodes 임

		// 하일라이트 처리
		// 		- edges 정보를 토대로 관련된 다른 노드들까지 하일라이트 처리
		//g.nodes, g.nodes.rect
		element
			.on("mouseover", (d:any, i:number, els:SVGElement[]) => {
				let el:SVGElement = isNodes?<any>els[i].parentElement: els[i];

				// 해당 element 와 관계된 element 조회 & CSS 적용
				let map:Map<string, SVGElement> = GraphUI.getRelationNodes(el, sourceEdge);
				map.forEach( (value:SVGElement, k:string) => {if(value) value.classList.add("focus")});

				// 커스텀 이벤트 발생
				if(options.events) {
					if(el.tagName=="line" && options.events.edge && options.events.edge.mouseover) options.events.edge.mouseover(el, d);
					else  if(options.events.node && options.events.node.mouseover) options.events.node.mouseover(el, d);
				}

				d3.event.stopPropagation(); //evnet bubble stop

			}).on("mouseout", (d:any, i:number, els:SVGElement[]) => {
				let el:SVGElement = isNodes?<any>els[i].parentElement: els[i];
				
				// 해당 element 와 관계된 element 조회 & CSS 적용
				let map:Map<string, SVGElement> = GraphUI.getRelationNodes(el, sourceEdge);
				map.forEach( (value:SVGElement, k:string) => {if(value) value.classList.remove("focus")});

				// 커스텀 이벤트 발생
				if(options.events) {
					if(el.tagName=="line" && options.events.edge && options.events.edge.mouseout) options.events.edge.mouseout(el);
					else  if(options.events.node && options.events.node.mouseout) options.events.node.mouseout(el);
				}


			});

	}

	/**
	 * 클릭시 선택 취소 
	 * 
	 * @param svg SVG element
	 * @param options  옵션
	 */
	public static onUnSelected(svg:d3Select.Selection<SVGSVGElement, any, SVGElement, any>, options:ConfigModel.Mesh) {

		svg.on("click", ()=> {
			svg.node().querySelectorAll(".selected").forEach((value:SVGElement) =>  value.classList.remove("selected"));
			if(options.events && options.events.unselected) options.events.unselected();	//커스텀 이벤트

		});

	}

	/**
	 * 클릭시 선택 
	 * 
	 * @param element 대상 element
	 * @param mapEdges 노드 연결정보 
	 * @param options  옵션
	 */
	public static onSelected(element:d3.Selection<SVGElement,any,SVGElement,any>, sourceEdge:model.EdgeSource, options:ConfigModel.Mesh) {

		if(!element || !element.node()) return;

		const isNodes:boolean = (element.node().tagName == "rect");	//rect nodes 임

		// 하일라이트 처리
		// 		- edges 정보를 토대로 관련된 다른 노드들까지 하일라이트 처리
		//g.nodes, g.nodes.rect
		element
			.on("click", (d:any, i:number, els:SVGSVGElement[]) => {
				let el:SVGElement = isNodes?<any>els[i].parentElement: els[i];

				// CSS 클래스 "selected" SVG에서 모두 제거
				el.ownerSVGElement.querySelectorAll(".selected").forEach((value:SVGElement) =>  value.classList.remove("selected"));

				// 해당 element 와 관계된 element 조회 & CSS 적용
				let map:Map<string, SVGElement> = GraphUI.getRelationNodes(el, sourceEdge);
				map.forEach( (value:SVGElement, k:string) => value.classList.add("selected"));

				// 커스텀 이벤트 발생
				//		UI 전달 파라메터
				//		- edge :  데이터
				//		- 선택된 엘리먼트의 class 가 "nodes"(박스) 이면 children 데이터 모두를 전달하지 않고 박스내 데이터(array) 값인 d.nodes 값만 넘겨준다.
				if(el.classList.contains("edge") && options.events.edge.selected) options.events.edge.selected(el, d);
				else  if(options.events.node.selected) options.events.node.selected(el, el.classList.contains("nodes") ? d.nodes: d);

				d3.event.stopPropagation(); //evnet bubble stop

			})

	}
	
}