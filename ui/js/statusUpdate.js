
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
printProgress=null;
jobStatusVisable=false;
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

function refreshStatus(){
	
	$.ajax({url:"/status"}).done(function(status,o){
		refreshUi(status);
		setTimeout(refreshStatus,10000);
	});
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

var chart;
$(document).ready(function(){
	printProgress=new ProgressBar("#printProg")
	//setTimeout(refreshStatus,10000);
	
	jobStatus(false);
	refreshStatus();
});