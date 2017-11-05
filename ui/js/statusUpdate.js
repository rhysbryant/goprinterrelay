
function ProgressBar(selecter){
	var progBar=$(selecter)
	var me=this
	var progress=0;
	this.val=function(val){
		if( typeof(val) == 'undefined' ){
			return progress;
		}
		else if( typeof(val) != 'number' ){
			throw 'expected Number'
		}
		progress=val;
		progBar
		.css({width:progress+"%"})
		.html(progress+"%");
	}
	this.hide=function(){
		progBar.hide();
	}
	this.show=function(){
		progBar.show();
	}
}
currentPage=null
printProgress=null;
currentPage=null
jobStatusVisable=false;
imageStreamer=null;
function refreshUi(data){
	
	$("#vElapsed").html(data.elapsedTime)
	$("#vEstmated").html(data.estmatedTime)
	$("#tExtruder1").html(data.temperature)
	$("#pStatus").html(data.status)
	$("#fTotal").html((data.filament.total/1000));
	$("#fRemaining").html(data.filament.remaining/1000);
	$("#pSerial").html(data.printerInfo.serial)
	$("#pModel").html(data.printerInfo.type)
	
	printProgress.val(data.printProgress);
	
	if(jobStatusVisable && data.estmatedTime==0 ){
		jobStatus(false)
		jobStatusVisable=false;
	}else if(!jobStatusVisable && data.estmatedTime>0){
		jobStatus(true)
		jobStatusVisable=true;		
	}
}

function refreshStatus(onlyOnce){
	
	$.ajax({url:"/status"}).done(function(status,o){
		refreshUi(status);
		
		if (typeof(onlyOnce)=='undefined' || !onlyOnce){
			setTimeout(refreshStatus,10000);
		}
	});
}

function onStatusUpdateReceived(evtData){
	var received_msg = evtData.data;
	if( typeof(received_msg) == "string"){
		received_msg.type="application/json";
		refreshUi(JSON.parse(received_msg));
	} 
		
}

function startStatusPush(){
	
	ws = new WebSocket("ws://"+location.host+"/statusPush");
	ws.onmessage=onStatusUpdateReceived;
	ws.onerror=refreshStatus;
}

function jobStatus(show){
	if(show === true){
		$("#printJobProgress").show();
		$("#noPrintJobNotice").hide();
	}else{
		$("#noPrintJobNotice").show();
		$("#printJobProgress").hide();
	}
}

function ajaxFormSubmit(){
	var form=$(this).parent()
	opt={
		target:form.find('.response')   
	}
	form.find(".form-response-card").show();
	form.find('form').ajaxSubmit(opt); 
}

function createFrormRow(label,type,name){
	var tp=$("#tools-form-item-template").clone()
	tp
	.find("label")
	.html(label)
	.attr('for','i'+name)
	tp
	.find("input")
	.attr('type',type)
	.attr('id','i'+name)
	.attr('name',name);
	tp
	.css({display:"block"});
	
	return tp;
}

function createForm(tool,id){
	var tp=$("#tools-formtemplate").clone();
	tp
	.removeAttr("id")
	.find(".form-name").html(tool.name)
	.show()
	form=tp.find("form")
	form.find("input[name=toolId]").val(id)
	for(var i in tool.formfields){
		arg=tool.formfields[i]
		formItem=createFrormRow(arg.name,arg.type,"toolArg"+i)
		form.append(formItem)
	}
	tp.find("button").click(ajaxFormSubmit)
	tp.find(".formcontainer").appendTo("#toolsPage")
}


function getTools(){
	$.ajax({"url":"tools"}).done(
	function(toolsObj){
		for(var item in toolsObj){
			var tool=toolsObj[item]
			console.log(tool);
			createForm(tool,item);
		}
		
	});
}

function showHideHelp(){
	var c=$(this);
	if(c.attr('showen') == true){
		$(c.attr('href')).hide();
		c.removeAttr('showen')
	}else{
		c.attr('showen',true);
		$(c.attr('href')).show();
	}
}

function changeTag(t){
	pageContainer=currentPage.attr('href')
	$(pageContainer).hide();
	currentPage.removeClass('active');
	
	currentPage=$(this)
	newPage=currentPage.addClass('active').attr('href')
	$(newPage).show();
	
}

$(document).ready(function(){
	printProgress=new ProgressBar("#printProg")
	//setTimeout(refreshStatus,10000);
	loadApplicationConfig();
	getTools();
	jobStatus(false);
	refreshStatus(true);
	startStatusPush();
	
	$("a.helpLink").click(showHideHelp)
	$("#pageNav li a").click(changeTag)
	currentPage=$("#pageNav li a.active")
	imageStreamer=new ImageStreamer($("#camGrab"),$("#filmstrip"),$("#streamStartStop"))
});

function loadApplicationConfig(){
	$.ajax({url:"/applicationInfo"}).done(function(appInfo,o){
		applySettings(appInfo);
	});
}

function applySettings(appInfo){
	$("#appVersion").html(appInfo.version);
	if (appInfo.featureConfig.camera.autoStart){
		imageStreamer.start();
	}
}