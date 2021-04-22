
// =====================================================================================================================
// Login Related Submit Actions
// ---------------------------------------------------------------------------------------------------------------------
//
// Templates:
// 		submitLogin 				| mt/LoginPage.html
//		submitValidate2faPin 		| mt/X2faPin.html
// 		submitRecoverPwPt1 			|
//		submitRecoverPasswordPt1    | mt/RecoverForgottenPassword.html
// 		submit2faSetupComplete  	| mt/DisplayRegistrationQR.html
// 		submitRegister 			    | mt/Register.html
//  	forgottenPassword 			| mt/LoginPage.html					| mt/Register.html
//
// Function coneect to form:
//  forgottenPassword 		mt/RecoverForgottenPassword.html
// ---------------------------------------------------------------------------------------------------------------------





// ---------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------------------------------
function submitLogin ( event ) {
	if ( event && event.preventDefault && typeof event.preventDefault == "function") { event.preventDefault(); }
	console.log ( "call to: submitLogin");
	var data = {
	   	  "username"		: $("#username").val()
		, "password"		: $("#password").val()
		, "__method__"		: "POST"
		, "_ran_" 			: ( Math.random() * 10000000 ) % 10000000
	};
	submitItData ( event, data, "/api/v2/login", function(data){
		console.log( "data=", data );
		user_id = data.user_id; // sample: -- see bottom of file: www/js/pdoc-form02.js
		g_user_id = data.user_id; // sample: -- see bottom of file: www/js/pdoc-form02.js
		g_username = $("#username").val();
		if ( data["real_name"] ) {
			$("#user_real_name").text(": "+data["real_name"]);
		}
		if ( data.Auth2faEnabled === "Yes" ) {
			renderX2faPin();
		} else {
			if ( window.localStorage ) {
				localStorage.setItem("jwt_token",data.jwt_token);
			}
			LoggInDone( data.jwt_token, true );
		}
	}, function(data) {
		console.log( "ERROR: ", data );
		renderError( "Failed to Login - Network communication failed.", "Failed to communicate with the server." );
	});
}

// ---------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------------------------------
function submitValidate2faPin ( event ) {
	if ( event && event.preventDefault && typeof event.preventDefault == "function") { event.preventDefault(); }
	console.log ( "call to: submitValidate2faPin");
	var data = {
	   	  "username"		: g_username
		, "pin2fa"			: $("#pin2fa").val()
		, "user_id"			: g_user_id
		, "__method__"		: "POST"
		, "_ran_" 			: ( Math.random() * 10000000 ) % 10000000
	};
	submitItData ( event, data, "/api/v2/2fa-validate-pin", function(data){
		console.log( "data=", data );
		LoggInDone( data.jwt_token, true );
	}, function(data) {
		console.log( "ERROR: ", data );
		renderError( "Failed to Login - Network communication failed.", "Failed to communicate with the server." );
		render5SecClearMessage() ;
	});
}

// function submitRecoverPwPt1 ( event ) {
// 	if ( event ) { event.preventDefault(); }
// 	console.log( "call to: submitRecoverPwPt1");
// 	renderRecoverForgottenPassword();
// }

// ---------------------------------------------------------------------------------------------------------------------
// recovery token form submit
// ---------------------------------------------------------------------------------------------------------------------
function submitRecoverPasswordPt1 ( event ) {
	if ( event && event.preventDefault && typeof event.preventDefault == "function") { event.preventDefault(); }
	console.log ( "call to: submitRecoverPasswordPt1");
	var data = {
	   	  "email"			: $("#email").val()
		, "__method__"		: "POST"
		, "_ran_" 			: ( Math.random() * 10000000 ) % 10000000
	};
	submitItData( event, data, "/api/v2/recover-password-pt1", function(data){
		console.log( "data=", data );
		// renderMessage ( "Successful Login", "You are now logged in<br>");
		// xyzzy - change form to show Enter Token - or put up new form.
	}, function(data) {
		console.log ( "ERROR: ", data );
		renderError ( "Failed to Login - Network communication failed.", "Failed to communicate with the server." );
	});
}





// ---------------------------------------------------------------------------------------------------------------------
// Called after registration when the QR code is displayed for setup of 2fa device.
// ---------------------------------------------------------------------------------------------------------------------
function submit2faSetupComplete ( event ) {
	if ( event && event.preventDefault && typeof event.preventDefault == "function") { event.preventDefault(); }
	console.log ( "call to: submit2faSetupComplete ");
	renderX2faPin ( event );
}



// ---------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------------------------------
function submitRegister ( event ) {
	if ( event && event.preventDefault && typeof event.preventDefault == "function") { event.preventDefault(); }
	console.log ( "call to: submitRegister ");
	var data = {
	   	  "username"			: $("#username").val()
		, "password"			: $("#password").val()
		, "real_name"			: $("#real_name").val()
		, "again"				: $("#password").val()
		, "email"				: $("#email").val()
		, "registration_token"	: $("#registration_token").val()
		, "__method__"			: "POST"
		, "_ran_" 				: ( Math.random() * 10000000 ) % 10000000
	};
	submitItData ( event, data, "/api/v2/register", function(data){
		console.log ( "data=", data );
		user_id = data.user_id; // sample: -- see bottom of file: www/js/pdoc-form02.js
		g_user_id = data.user_id; // sample: -- see bottom of file: www/js/pdoc-form02.js
		if(data.Auth2faEnabled === "yes" || data.Auth2faEnabled === "Yes") {
			renderDisplayRegistrationQR( data );	// qr_url
		} else {
			if ( window.localStorage ) {
				localStorage.setItem("jwt_token",data.jwt_token);
			}
			renderMessage ( "Successful Registration", "You are now logged in<br>");
			LoggInDone( data.jwt_token, true );
			if ( window["postRegistrationSetup"] ) {
				// var fx = window["postRegistrationSetup"];
				// fx ( data );
				callFunc ( "postRegistrationSetup", data );
			}
		}
	}, function(data) {
		console.log ( "ERROR: ", data );
		renderError ( "Failed to Registration - Network communication failed.", "Failed to communicate with the server." );
	});
}


// ---------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------------------------------
function forgottenPassword ( event ) {
	if ( event && event.preventDefault && typeof event.preventDefault == "function") { event.preventDefault(); }
	console.log ( "call to: forgottenPassword ");
	renderRecoverForgottenPassword();
}

