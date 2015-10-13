(function () {
	angular
		.module('mainApp', [
			'ngRoute'
		])
		.config(['$routeProvider', function ($routeProvider) {
			$routeProvider
				.when('/', {
					templateUrl: 'partials/genres',
					controller: 'GenreCtrl'
				})
				.otherwise({
					redirectTo: '/'
				})
		}])
})();