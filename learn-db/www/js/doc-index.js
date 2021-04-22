
// glboal constant applicaiton name
window.g_data = {
	  "app_name": "Demo Auth"
	, "title": ""
	, "msg": ""
	, "title": ""
};

var jwt_token = "" ;
var xsrf_token = "" ;
var g_user_id;
var g_token;
var g_username = "";
var x; 		// temp to debug stuff 3
var g_qr_url = "";   // URL of QR Code for registration setup of 2fa

var LoggInDone;	// call function on successful login.
var LoggOut;	// call function on logout clicked.

var g_item_id ; // From form33 - upload
var g_event_id ; // From form33 - upload
var g_app_name = 'www.pure-imagination-server.com';

var g_page_state = "";	// not used yet

function tss(tmpl, ns, fx) {
	if ( typeof fx == "undefined" ) { fx = {} };
	var p1 = tmpl.replace(/%{([A-Za-z0-9_|.,]*)%}/g, function(j, t) {
			var pl;
			var s = "";
			var a = t.split("|");
			for ( var i = 0; i < a.length; i++ ) {
// console.log ( "ts1: a["+i+"] =["+a[i]+"]" );
				pl = a[i].split(",");
// console.log ( "ts1: pl[0] =["+pl[0]+"]"+"  typeof for this =="+typeof fx[pl[0]] );
				if ( typeof fx [ pl[0] ] == "function" ) {
					var f = fx[ pl[0] ];
					s = f(s,ns,pl);									// Now call each function as a pass thru-filter
				} else if ( typeof ns [ pl[0] ] == "string" || typeof ns [ pl[0] ] == "number" ) {
					s = ns[ pl[0] ];
				} else {
					s = "";
				}
			}
			return s;
		});
	return p1;
};

function renderError ( title, msg, data ) {
	if ( data ) {
		title = tss(title, data);
		msg = tss(msg, data);
	}
	var form = [ ''
		,'<div style="margin-bottom:24px;">'
			,'<div class="row">'
				,'<div class="col-sm-12">'
					,'<div class="card bg-danger">'
						,'<div class="card-header"><h4 style="color:white;">'+title+'</h4></div>'
						,'<div class="card-body bg-light">'
							,'<div>'+msg+'</div>'
						,'</div>'
					,'</div>'
				,'</div>'
			,'</div>'
		,'</div>'
	].join("\n");
	$("#msg").html(form);
}

function renderMessage ( title, msg, data ) {
	if ( data ) {
		title = tss(title, data);
		msg = tss(msg, data);
	}
	var form = [ ''
		,'<div style="margin-bottom:24px;">'
			,'<div class="row">'
				,'<div class="col-sm-12">'
					,'<div class="card bg-success">'
						,'<div class="card-header"><h4 style="color:white;">'+title+'</h4></div>'
						,'<div class="card-body bg-light">'
							,'<div>'+msg+'</div>'
						,'</div>'
					,'</div>'
				,'</div>'
			,'</div>'
		,'</div>'
	].join("\n");
	$("#msg").html(form);
}

function renderClearMessage() {
	$("#msg").html("");
}
function render5SecClearMessage() {
	setTimeout(function() {
		$("#msg").html("");
	},5000);
}

// xyzzy - need to ad a title for errors
function submitItData ( event, data, action, succ, erro ) {
	if ( event && event.preventDefault && typeof event.preventDefault == "function") { event.preventDefault(); }
	console.log ( "form data (passed): ", data);
	$.ajax({
		type: 'GET',
		url: action,
		data: data,
		success: function (data) {
			console.log ( "success AJAX", data );
			$("#output").text( JSON.stringify(data, null, 4) );
			if ( succ ) {
				succ(data);
			}
		},
		error: function(resp) {
			$("#output").text( "Error!"+JSON.stringify(resp) );
			// alert("got error status="+resp.status+" "+resp.statusText);
			if ( erro ) {
				erro(resp);
			} else {
				var msg = resp.statusText;
				renderError ( "xyzzy-title:doc-index.js:126", msg );
				render5SecClearMessage ( );
				if ( msg === "Unauthorized" ) {
					LoggOut();
					renderLoginPage() ;
				}
			}
		}
	});
}

function fetchRefData ( url, succ, erro ) {
	var data = { "auth_key": g_auth_key, "_ran_": ( Math.random() * 10000000 ) % 10000000 };
	$.ajax({
		type: 'GET',
		url: url,
		data: data,
		success: function (data) {
			$("#output").text( JSON.stringify(data, null, 4) );
			if ( succ ) {
				succ(data.data);
			}
		},
		error: function(resp) {
			$("#output").text( "Error!"+JSON.stringify(resp) );
			// alert("got error status="+resp.status+" "+resp.statusText);
			if ( erro ) {
				erro(resp);
			}
		}
	});
}

// Function will return "true" if the user is logged in.
// This is based on a non-zero length jwt_token in local storage.
function isLoggedIn() {
	var rv = false;
	if ( window.localStorage ) {
		var x = localStorage.jwt_token;
		if ( x && x.length > 0 ) {
			jwt_token = x;
			rv = true;
		}
	}
	return ( rv );
}

// Same as isLoggedIn - but sets UI to be in logged in mode.
// Called only at startup time.
function isLoggedInStartupCheck() {
	var rv = false;
	if ( window.localStorage ) {
		var x = localStorage.jwt_token;
		if ( x && x.length > 0 && x === "undefined" ) {
			window.localStorage.removeItem ( "jwt_token" );
			// window.localStorage.removeStorage ( "jwt_token" );
			x = undefined;
			console.log ( "Removed a 'undefined' from local storage" );
		}
		if ( x && x.length > 0 && x !== "undefined" ) {
			jwt_token = x;
			LoggInDone (jwt_token);
			rv = true;
			console.log ( " " );
			console.log ( "   You Are Logged In " );
			console.log ( " " );
		}
	}
	return ( rv );
}

function doLogin ( event ) {
	if ( event ) { event.preventDefault(); }
	console.log ( "doLogin");
	renderLoginPage() ;
}
function doRegister ( event ) {
	if ( event ) { event.preventDefault(); }
	console.log ( "doRegister");
	// renderRegistration();
	renderRegister();
}

LoggInDone = function ( JWTToken, showMsg ) {
	console.log ( "LoggInDone");
	$("#body").html("<span></span>");
	jwt_token = JWTToken;
	if ( window.localStorage ) {
		localStorage.setItem("jwt_token",jwt_token);
	}
	SetupJWTBerrer();
	$("#login").html("Logout");
	$("#login").click(LoggOut);
	$(".show-anon").hide();
	$(".show-logged-in").show();
	if ( showMsg ) {
		renderMessage( "Successful Login", "You are now logged in<br>");
		render5SecClearMessage() ;
	}
}

LoggOut = function ( event, undefined ) {
	console.log ( "LoggOut");
	if ( event ) {
		event.preventDefault();
	}
	xsrf_token = "" ;
	jwt_token = "" ;
	localStorage.setItem("jwt_token",undefined);
	window.localStorage.removeItem ( "jwt_token" );
	//	window.localStorage.removeStorage ( "jwt_token" );
	//	var x = localStorage.jwt_token;
	//	if ( x && x.length > 0 && x === "undefined" ) {
	//		window.localStorage.removeStorage ( "jwt_token" );
	//		x = undefined;
	//		console.log ( "Removed a 'undefined' from local storage' );
	//	}
	$("#login").html("Login");
	$("#login").click(doLogin);
	$(".show-anon").show();
	$(".show-logged-in").hide();
	renderMessage ( "Logged Out", "You are now logged out.<br>");
	render5SecClearMessage();
	renderTitleMainPage(); // render Login Welcome Message
}





$("#login").click(doLogin);
$("#register").click(doRegister);
$(".show-logged-in").hide();
if ( isLoggedInStartupCheck() ) {
	console.log ( "jwt_token - page refresh -(on next line)- ", jwt_token );
}
$("#logout").click(LoggOut);







function submitIt ( event, data, action, succ, erro ) {
	if ( event && event.preventDefault && typeof event.preventDefault == "function") { event.preventDefault(); }

	// xyzzy - add in _ran_ to data

	$.ajax({
		type: 'GET',
		url: action,
		data: data,
		success: function (data) {
			if ( succ ) {
				succ(data);
			}
			$("#output").text( JSON.stringify(data, null, 4) );
		},
		error: function(resp) {
			$("#output").text( "Error!"+JSON.stringify(resp) );
			if ( erro ) {
				erro(data);
			}
			// alert("got error status="+resp.status+" "+resp.statusText);
		}
	});
}

$("#getStatus").click(function(event){
	submitIt ( event, {}, "/api/v1/status",
		function(data) {
			$("#body").html( "<pre>"+JSON.stringify(data)+"</pre>" );
		}
	);
});

function renderStatus ( event )  {
	if ( event && event.preventDefault && typeof event.preventDefault == "function") { event.preventDefault(); }
	submitIt ( event, {}, "/api/v1/status",
		function(data) {
			$("#body").html( "<pre>"+JSON.stringify(data)+"</pre>" );
		}
	);
}


// ------------------------------------------------------------------------------------------------------------------------------------
//					<a class="dropdown-item" href="#" id="auth-change-password" > Change Password</a>			
// ------------------------------------------------------------------------------------------------------------------------------------

function submitChangePassword(event) {
	if ( event && event.preventDefault && typeof event.preventDefault == "function") { event.preventDefault(); }
	console.log ( "call to: submitChangePassword()" );
//	var data = {
//	   	  "old_pw"			: $("#old_password").val()
//		, "new_password"	: $("#password").val()
//		, "confirm"			: $("#again").val()
//		, "__method__"		: "POST"
//		, "_ran_" 			: ( Math.random() * 10000000 ) % 10000000
//	};
	$("#_ran_").val( ( Math.random() * 10000000 ) % 10000000 );
	var frm = $( "#form07" );
	var data = frm.serialize();
	var url = frm.attr('action');
	submitItData ( event, data, url, function(data){
		console.log( "data=", data );
		if ( data.jwt_token ) {
			LoggInDone( data.jwt_token, false );
			renderMessage( "Password Changed", "You have successfully chagned your password.<br>");
			render5SecClearMessage() ;
		}
	}, function(data) {
		console.log( "ERROR: ", data );
		renderError( "Failed to Login - Network communication failed.", "Failed to communicate with the server." );
	});
}





// ------------------------------------------------------------------------------------------------------------------------------------
//					<a class="dropdown-item" href="#" id="auth-new-dev-account" > New Device Account</a>	
// ------------------------------------------------------------------------------------------------------------------------------------

function submitDevUnPwCreate(event) {
	if ( event && event.preventDefault && typeof event.preventDefault == "function") { event.preventDefault(); }
	console.log ( "call to: submitDevUnPwCreate()" );
	$("#_ran_").val( ( Math.random() * 10000000 ) % 10000000 );
	var frm = $( "#form08" );
	var data = frm.serialize();
	var url = frm.attr('action');
	console.log ( "data=", data, 'url=', url );
	submitItData ( event, data, url, function(data){
		console.log( "data=", data );
		renderMessage( "Device Account Created", "You have successfully created an account for a device to login with." );
		render5SecClearMessage() ;
	}, function(data) {
		console.log( "ERROR: ", data );
		renderError( "Failed to Create Device Account", "Unable to Create Account:<br>%{msg%}<br>errorid=%{id%}", data );
	});
}



// ------------------------------------------------------------------------------------------------------------------------------------
//					<a class="dropdown-item" href="#" id="admin-list-accounts" > List Accounts </a>			
// ------------------------------------------------------------------------------------------------------------------------------------

// xyzzy82

function submitListAccounts(event) {
	if ( event && event.preventDefault && typeof event.preventDefault == "function") { event.preventDefault(); }
	console.log ( "call to: submitListAccounts()" );
	var url = "/api/v2/list-sub-acct";
	var rv = ( Math.random() * 10000000 ) % 10000000;
	submitItData ( event, { "_ran_": rv }, url, function(data){
		console.log( "data=", data );
		// renderMessage( "Device Account Created", "You have successfully created an account for a device to login with." );
		// render5SecClearMessage() ;
	}, function(data) {
		console.log( "ERROR: ", data );
		renderError( "Failed to Create Device Account", "Unable to Create Account:<br>%{msg%}<br>errorid=%{id%}", data );
	});
}



















// ------------------------------------------------------------------------------------------------------------------------------------
// JWT Token Setup for AJAX calls
// ------------------------------------------------------------------------------------------------------------------------------------
function SetupJWTBerrer() {
	$.ajaxSetup({
		beforeSend: function(xhr) {
			if ( jwt_token && jwt_token !== "" ) {
				xhr.setRequestHeader('Authorization', 'bearer '+jwt_token);
			}
			if ( xsrf_token != "" ) {
				xhr.setRequestHeader('X-Xsrf-Token', xsrf_token);
			}
		}
		,dataFilter: function(data, type) {
			var prefix = ['//', ')]}\'', 'while(1);', 'while(true);', 'for(;;);'], i, l, pos;

			// console.log ( "dataFilter: data type", type );

			if (type && type != 'json' && type != 'jsonp') {
				return data;
			}

			// console.log ( "dataFilter: raw data before remove of prefix.", data );

			var dl = data.length;	 // data length 
			for (i = 0, l = prefix.length; i < l; i++) {
				var pl = prefix[i].length; // pattern lenght
				// console.log ( "dataFilter: raw substr -={" + data.substring(0,pl) + "}=-" );
				if ( dl >= pl && data.substring(0,pl) === prefix[i] ) {
					return data.substring(pl);
				}
			}

			// console.log ( "dataFilter: data after remove of prefix.", data );

			return data;
		}
	});
}

// ------------------------------------------------------------------------------------------------------------------------------------
// Ajax Upload Setup
// ------------------------------------------------------------------------------------------------------------------------------------

// Function that will allow us to know if Ajax uploads are supported
function supportAjaxUploadWithProgress() {
	return supportFileAPI() && supportAjaxUploadProgressEvents() && supportFormData();

	// Is the File API supported?
	function supportFileAPI() {
		var fi = document.createElement('INPUT');
		fi.type = 'file';
		return 'files' in fi;
	};

	// Are progress events supported?
	function supportAjaxUploadProgressEvents() {
		var xhr = new XMLHttpRequest();
		return !! (xhr && ('upload' in xhr) && ('onprogress' in xhr.upload));
	};

	// Is FormData supported?
	function supportFormData() {
		return !! window.FormData;
	}
}

function displayFileUploadSupported() {

	// Actually confirm support
	if (supportAjaxUploadWithProgress()) {
		// Ajax uploads are supported!
		// Change the support message and enable the upload button
		var notice = document.getElementById('support-notice');
		var uploadBtn = document.getElementById('upload-button-id');
		notice.innerHTML = "Your browser supports HTML uploads. Go try me! :-)";
		uploadBtn.removeAttribute('disabled');

		// Init the Ajax form submission
		initFullFormAjaxUpload();

		// Init the single-field file upload
		initFileOnlyAjaxUpload();
	}

}

function initFullFormAjaxUpload() {
	var form = document.getElementById('form-id');
	form.onsubmit = function() {
		// FormData receives the whole form
		var formData = new FormData(form);

		// We send the data where the form wanted
		var action = form.getAttribute('action');

		// Code common to both variants
		sendXHRequest(formData, action);

		// Avoid normal form submission
		return false;
	}
}

function UploadTheFile(ty) {
	var formData = new FormData();

	// Since this is the file only, we send it to a specific location
	var action = '/upload';

	// FormData only has the file
	var fileInput = document.getElementById('file-id');
	var file = fileInput.files[0];
	formData.append('file', file);
	formData.append('id', g_item_id);
	formData.append('type', ty);
	formData.append("app", g_app_name);

	// Code common to both variants
	sendXHRequest(formData, action);
}

function initFileOnlyAjaxUpload() {
	var uploadBtn = document.getElementById('upload-button-id');
	uploadBtn.onclick = function (evt) {
		var formData = new FormData();

		// Since this is the file only, we send it to a specific location
		var action = '/upload';

		// FormData only has the file
		var fileInput = document.getElementById('file-id');
		var file = fileInput.files[0];
		formData.append('file', file);
		formData.append('id', "15e42502-e7a5-44e2-6920-b410b9308412");

		// Code common to both variants
		sendXHRequest(formData, action);
	}
}

// Once the FormData instance is ready and we know
// where to send the data, the code is the same
// for both variants of this technique
function sendXHRequest(formData, uri) {
	// Get an XMLHttpRequest instance
	var xhr = new XMLHttpRequest();

	// Set up events
	xhr.upload.addEventListener('loadstart', onloadstartHandler, false);
	xhr.upload.addEventListener('progress', onprogressHandler, false);
	xhr.upload.addEventListener('load', onloadHandler, false);
	xhr.addEventListener('readystatechange', onreadystatechangeHandler, false);

	// Set up request
	xhr.open('POST', uri, true);

	// Fire!
	xhr.send(formData);
}

// Handle the start of the transmission
function onloadstartHandler(evt) {
	var div = document.getElementById('upload-status');
	div.innerHTML = 'Upload started.';
}

// Handle the end of the transmission
function onloadHandler(evt) {
	var div = document.getElementById('upload-status');
	div.innerHTML += '<' + 'br>File uploaded. Waiting for response.';
}

// Handle the progress
function onprogressHandler(evt) {
	var div = document.getElementById('progress');
	var percent = evt.loaded/evt.total*100;
	div.innerHTML = 'Progress: ' + percent + '%';
}

// Handle the response from the server
function onreadystatechangeHandler(evt) {
	var status, text, readyState;

	try {
		readyState = evt.target.readyState;
		text = evt.target.responseText;
		status = evt.target.status;
	}
	catch(e) {
		return;
	}

	if (readyState == 4 && status == '200' && evt.target.responseText) {
		var status = document.getElementById('upload-status');
		status.innerHTML += '<' + 'br>Success!';

		var result = document.getElementById('result');
		result.innerHTML = '<p>The server saw it as:</p><pre>' + evt.target.responseText + '</pre>';
	}
}

