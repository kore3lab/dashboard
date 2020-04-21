import * as d3			from "d3";
import * as d3Select	from "d3-selection";
import {Tree as model}	from "../../model/graph.model";
import {ConfigModel}	from "../../model/models";
import {Transform}		from "../../utils/transform";
import {GraphBase}		from "../graph.base";


require("./graph.css");
/**
 * Topology 그래프 랜더러
 */
export class RbacGraph extends GraphBase {


	/**
	 * (abstract) 랜더링
	 * @param data 토플로지를 위한 k8s 데이터 (model.K8s)
	 */
	public populate(conf:ConfigModel.Config, svgEl:d3Select.Selection<SVGSVGElement, any, SVGElement, any>, bounds:DOMRect, outlineEl:d3Select.Selection<SVGGElement,any,SVGElement,any>) {
		
		if(!conf.data) return;

		// svg > defs
		if(svgEl.select("defs").size() == 0) svgEl.append("defs").call(RbacGraph.renderDefs, conf);

		// data 가공
		let data:d3.HierarchyNode<model.Tree> = d3.hierarchy(conf.data, d =>d.children);
		d3.tree().nodeSize([15,15])(data);

		// body
		let bodyEl:d3Select.Selection<SVGGElement,any,SVGElement,any> = outlineEl.append("g");	
		bodyEl.call(RbacGraph.renderRoot, data);


		// 마지막 노드를 제외한 노드들중 최대 높이 기준으로 그리기 옵션
		let max:number = 0
		bodyEl.selectAll(".node").each((d:any, i:number, m:SVGElement[]) => {
			console.log(d);
			if(d.children) max = Math.max(max, m[i].getBoundingClientRect().width);
		});
		d3.tree().nodeSize([15,max])(data);
		bodyEl.selectAll(".node").remove();
		bodyEl.selectAll(".link").remove();
		bodyEl.call(RbacGraph.renderRoot, data);

		const rect:any = bodyEl.node().getBoundingClientRect();
		console.log(rect);
		Transform.instance(bodyEl.node())
			.translate(rect.x<0?rect.x*-1:0, 22);


		// d3.tree().nodeSize([max, 100])(data);
		// bodyEl.call(RbacGraph.renderRoot, data);
		// console.log(max);

		// // 사이즈에 꽉차게 ratio 조정
		// const body:SVGElement = bodyEl.node();
		// const rect:any = body.getBoundingClientRect();
		// const ratio = bounds.width / rect.width;

		// let trans:Transform = Transform.instance(body);
		// trans.translate(rect.x<0?rect.x*ratio*-1:0, rect.y<0?rect.y*ratio*-1:0).ratioScale(ratio);
		

	}

	private static renderRoot(parentEl:d3Select.Selection<SVGGElement,any,SVGElement,any>, data:d3.HierarchyNode<model.Tree>) {

		const iconR:number = 10

		// adds the links between the nodes
		parentEl.selectAll(".link")
			.data( data.descendants().slice(1) )
		.enter().append("path")
			.attr("class", "link")
			.attr("d", (d:any) => `M${d.y},${d.x} C${(d.y + d.parent.y)/2},${d.x} ${(d.y + d.parent.y) / 2},${d.parent.x} ${d.parent.y},${d.parent.x}`);

		// adds each node as a group
		let nodes:d3Select.Selection<SVGGElement,any,SVGElement,any> = parentEl.selectAll(".node")
			.data(data.descendants())
		.enter().append("g")
    		.attr("class", (d:any) => "node" + (d.children ? " node--internal" : " node--leaf"))
			.attr("transform", (d:any) => `translate(${d.y},${d.x})`);


		nodes
			.append("use")
			.attr("height", 24).attr("width", 24)
			.attr("xlink:href", (d:any) => `#ac_ic_node_${d.data.kind}`)
			.attr("x", -12)
			.attr("y", -12) 

		// adds the text to the node
		nodes.append("text")
		.attr("x", (d:any) => (d.children ? -(iconR+3) : iconR+3) )
		.style("text-anchor", (d:any) => (d.children ? "end" : "start"))
		.text((d:any) => d.data.name)
		.attr("y", (d:any, i:number, els:SVGTextElement[])=> (iconR*2 - els[i].getBoundingClientRect().height)*-1)

	}

	/**
	 * defs 정의
	 * 
	 * @param defsEl def 엘리먼트
	 */
	private static renderDefs(defsEl:d3.Selection<SVGDefsElement, any, SVGElement, any>) {

		defsEl.append("symbol").attr("id", "ac_ic_node_rolebinding").attr("viewBox", "0 0 24 24").append("path").attr("d", "M3.9 12c0-1.71 1.39-3.1 3.1-3.1h4V7H7c-2.76 0-5 2.24-5 5s2.24 5 5 5h4v-1.9H7c-1.71 0-3.1-1.39-3.1-3.1zM8 13h8v-2H8v2zm9-6h-4v1.9h4c1.71 0 3.1 1.39 3.1 3.1s-1.39 3.1-3.1 3.1h-4V17h4c2.76 0 5-2.24 5-5s-2.24-5-5-5z")
		defsEl.append("symbol").attr("id", "ac_ic_node_clusterrolebinding").attr("viewBox", "0 0 24 24").append("path").attr("d", "M3.9 12c0-1.71 1.39-3.1 3.1-3.1h4V7H7c-2.76 0-5 2.24-5 5s2.24 5 5 5h4v-1.9H7c-1.71 0-3.1-1.39-3.1-3.1zM8 13h8v-2H8v2zm9-6h-4v1.9h4c1.71 0 3.1 1.39 3.1 3.1s-1.39 3.1-3.1 3.1h-4V17h4c2.76 0 5-2.24 5-5s-2.24-5-5-5z")
		defsEl.append("symbol").attr("id", "ac_ic_node_role").attr("viewBox", "0 0 24 24").append("path").attr("d", "M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z")
		defsEl.append("symbol").attr("id", "ac_ic_node_clusterrole").attr("viewBox", "0 0 24 24").append("path").attr("d", "M11.99 2c-5.52 0-10 4.48-10 10s4.48 10 10 10 10-4.48 10-10-4.48-10-10-10zm3.61 6.34c1.07 0 1.93.86 1.93 1.93 0 1.07-.86 1.93-1.93 1.93-1.07 0-1.93-.86-1.93-1.93-.01-1.07.86-1.93 1.93-1.93zm-6-1.58c1.3 0 2.36 1.06 2.36 2.36 0 1.3-1.06 2.36-2.36 2.36s-2.36-1.06-2.36-2.36c0-1.31 1.05-2.36 2.36-2.36zm0 9.13v3.75c-2.4-.75-4.3-2.6-5.14-4.96 1.05-1.12 3.67-1.69 5.14-1.69.53 0 1.2.08 1.9.22-1.64.87-1.9 2.02-1.9 2.68zM11.99 20c-.27 0-.53-.01-.79-.04v-4.07c0-1.42 2.94-2.13 4.4-2.13 1.07 0 2.92.39 3.84 1.15-1.17 2.97-4.06 5.09-7.45 5.09z")
		defsEl.append("symbol").attr("id", "ac_ic_node_group").attr("viewBox", "0 0 24 24").append("path").attr("d", "M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z")
		defsEl.append("symbol").attr("id", "ac_ic_node_user").attr("viewBox", "0 0 24 24").append("path").attr("d", "M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z")


	}
	
	
};	
