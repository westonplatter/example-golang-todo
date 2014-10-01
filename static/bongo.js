
// bongo app - test app using backbone + golang
// basic todo app with categories

$(function(){

	// Backbone.emulateJSON = true;

	var post = Backbone.Model.extend({

		idAttribute: 'Id',

		// default values
		defaults: function() {
			return {
				Title: "default title",
				Category: "Work",
				State: "active",
				Dt_completed: 0,
				Dt_created: new Date(),
			};
		},

		// on initialize set defaults
		initialize: function() {
			if (!this.get("Title")) { this.set({"Title": this.defaults().Title}); }
			if (!this.get("State")) { this.set({"State": this.defaults().State}); }
			if (!this.get("Dt_completed")) { this.set({"Dt_completed": this.defaults().Dt_completed}); }
			if (!this.get("Dt_created")) { this.set({"Dt_created": this.defaults().Dt_created}); }
			if (!this.get("Category")) { this.set({"Category": this.defaults().Category}); }
		},

		urlRoot: "/todos/",

	});

	//------------------------------------------------

	var postList = Backbone.Collection.extend({
		model: post,
		url: "/todos/",
	});

	var posts = new postList;

	//------------------------------------------------

	var postView = Backbone.View.extend({

		tagName:  "div",

		// cache the template function
		template: _.template($('#item-template').html()),

		events: {
			"click .edit-icon"		: "edit_task",
			"blur .edit"			: "edit_close",
			"click .archive"        : "archive",
			"click .check-box"      : "complete_task",
		},

		initialize: function() {
			this.listenTo(this.model, 'change', this.render);
			this.listenTo(this.model, 'destroy', this.remove);
		},

		render: function() {
			var context = this.model.toJSON();
			var dt = new Date(context['Dt_created']);
			if (typeof(new Date().toLocaleFormat) == "undefined" ) {
				context['dt_created_fmt'] = dt.toLocaleString();
			}
			else {
				context['dt_created_fmt'] = dt.toLocaleFormat('%B %d, %Y %I:%M%p');
			}

			context['completed'] = ( context['Dt_completed'] > 0 ) ? "checked" : "";


			// parse for todo list
			var html = this.template(context);
			this.$el.html(html);
			return this;
		},


		edit_task: function(e) {
			this.$el.addClass("editing");
			this.$el.find(".edit-field").focus();

		},

		edit_close: function(e) {
			var el = this.$el.find(".edit-field");
			if (! el.val()) return;
			this.model.save({Title: el.val()});
			this.$el.removeClass("editing");
		},


		complete_task: function(e) {
			completed = this.model.get("Dt_completed");
			cval = ( completed > 0 ) ? 0 : Date.now();
			this.model.save({Dt_completed: cval, Id:this.model.get("Id")});

		},


		// destroy item, remove model from collection
		archive: function() {
			this.model.set({State: "archived", Id: this.model.get("Id")});
			this.model.destroy();
		},


	});


	//------------------------------------------------
	var AppView = Backbone.View.extend({

		el: $("#bongo"),

		events: {
			"click #save-task"				: "createOnSave",
			"keypress #new-task-title"		: "checkForEnter",
			"click #filter-btn li"			: "filter_category"
		},


		initialize: function() {

			this.title    = this.$("#new-task-title");
			this.category = this.$("#new-task-category");

			this.listenTo(posts, 'add', this.addOne);
			this.listenTo(posts, 'reset', this.addAll);
			this.listenTo(posts, 'all', this.render);

			this.main = $('#main');
			posts.fetch();
		},


		addOne: function(post) {
			var view = new postView({model: post});
			this.$( '#task-list' ).append( view.render().el );
		},


		addAll: function() {
			posts.each(this.addOne, this);
		},


		createOnSave: function(e) {
			if (!this.title.val()) return;

			posts.create({
				Title: this.title.val().trim(),
				Category: this.category.val().trim(),
				Dt_completed: 0,
				Dt_created: parseInt(new Date().getTime())
			});
			this.title.val('');			// reset form data
		},

		checkForEnter: function(e) {
			if (e.keyCode == 13) {
				this.createOnSave()
			}
		},


		filter_category: function(e) {	/* filter category list */

			//TODO: convert this mess to be data driven template
			//	    let backbone do the work

			var el = $(e.target);
			// normalize where click is coming from, could be icon, text or li
			// make all clicks go up to the LI tag
			if ( el.prop("tagName") != "LI" ) {
				el = el.parent();
			}

			var pel = el.parent();
			var gpel = pel.parent();
			category = el.find("span").html()
			category = category.replace( /\W/, '' );

			/* switch out active filter */
			active_elem = gpel.find("span.filter");
			active_elem.removeClass();
			active_elem.addClass("filter");
			cat_icon = category.toLowerCase() + "-icon";
			active_elem.addClass(cat_icon);
			active_elem.html(category);

			if ( category == "All" ) {
				$('.view').show();
			}
			else {
				$('.view').hide();
				$('.'+category+'-view').show();
				this.category.val(category);	// set new task category
			}
		}

	});

	// lets get the party started
	var App = new AppView;

});

$(function() {
	$('.close-btn').click(function() {
		window.location.hash = '';
	});

	$('#add-new-task-btn').click(function() {
		if (window.location.hash == '#enter-form') {
			window.location.hash = '';
		}
		else {
			window.location.hash='enter-form';
		}
	});
});

