from flask import Flask, render_template, request, redirect, url_for, flash
from flask_sqlalchemy import SQLAlchemy
from datetime import datetime, timedelta
from flask_apscheduler import APScheduler
from flask_login import current_user
from flask_login import login_required
import requests

def fetch_netflix_subscriptions(api_key, user_id):
    url = f"https://api.netflix.com/subscriptions?user_id={user_id}"
    headers = {"Authorization": f"Bearer {api_key}"}
    
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        netflix_subscriptions = response.json()
        return netflix_subscriptions
    else:
        print(f"Error fetching Netflix subscriptions. Status code: {response.status_code}")
        return None

# Usage
netflix_api_key = "your_netflix_api_key"
netflix_user_id = "user123"
netflix_data = fetch_netflix_subscriptions(netflix_api_key, netflix_user_id)
if netflix_data:
    print("Netflix Subscriptions:")
    print(netflix_data)


app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///subscriptions.db'
app.config['SECRET_KEY'] = 'secret_key'
app.config['SCHEDULER_API_ENABLED'] = True
db = SQLAlchemy(app)
scheduler = APScheduler()

class User(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(50), unique=True, nullable=False)
    email = db.Column(db.String(120), unique=True, nullable=False)
    password = db.Column(db.String(60), nullable=False)
    subscriptions = db.relationship('Subscription', backref='user', lazy=True)

class Subscription(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(100), nullable=False)
    cost = db.Column(db.Float, nullable=False)
    renewal_date = db.Column(db.Date, nullable=False)
    billing_cycle = db.Column(db.String(20), nullable=True)
    category = db.Column(db.String(50), nullable=True)
    user_id = db.Column(db.Integer, db.ForeignKey('user.id'), nullable=False)

db.create_all()

@app.route('/')
@login_required
def index():
    # Implement user authentication and get user-specific subscriptions
    user_subscriptions = Subscription.query.filter_by(user_id=current_user.id).all()
    return render_template('index.html', subscriptions=user_subscriptions)

# User Registration and Profile
@app.route('/register', methods=['GET', 'POST'])
def register():
    if request.method == 'POST':
        # Handle user registration form submission
        username = request.form['username']
        email = request.form['email']
        password = request.form['password']

        new_user = User(username=username, email=email, password=password)
        db.session.add(new_user)
        db.session.commit()

        flash('Account created successfully', 'success')
        return redirect(url_for('login'))

    return render_template('register.html')

def get_service_providers():
    # Implement logic to fetch service providers from a database or external API
    service_providers = ["Netflix", "Hotstar", "Electricity", "Milk Vendor"]
    return service_providers

# Subscription Dashboard
@app.route('/dashboard')
@login_required
def dashboard():
    user_subscriptions = Subscription.query.filter_by(user_id=current_user.id).all()
    service_providers = get_service_providers()
    return render_template('dashboard.html', subscriptions=user_subscriptions, service_providers=service_providers)

def process_payment(subscription):
     # Implement payment processing logic using a payment gateway API
    # Example: Stripe integration
    # stripe.api_key = "your_stripe_secret_key"
    # charge = stripe.Charge.create(amount=int(subscription.cost * 100), currency="usd", source="tok_visa")
    pass

def track_expense(subscription):
    # Implement logic to track expenses, store in the database, etc.
    pass

@app.route('/add_subscription', methods=['POST'])
def add_subscription():
    if request.method == 'POST':
        name = request.form.get['name']
        cost = float(request.form.get['cost', 0.0])
        renewal_date = datetime.strptime(request.form.get['renewal_date'], '%Y-%m-%d').date()
        billing_cycle = request.form.get['billing_cycle']
        category = request.form.get['category']

        new_subscription = Subscription(
            name=name,
            cost=cost,
            renewal_date=renewal_date,
            billing_cycle=billing_cycle,
            category=category,
            user_id=current_user.id
        )

        db.session.add(new_subscription)
        db.session.commit()
        process_payment(new_subscription)
        track_expense(new_subscription)
        flash('Subscription added successfully', 'success')

    return redirect(url_for('dashboard'))

# Automated Reminders
def send_reminder(subscription):
    # Implement sending reminders (e.g., through email or push notifications)
    pass

def send_notification(subscription):
    # Implement logic to send notifications (e.g., email, push notifications)
    pass

# Schedule reminders for upcoming renewals
@scheduler.task('interval', id='check_renewals', seconds=86400)  # Daily check
def check_renewals():
    subscriptions = Subscription.query.filter_by(user_id=current_user.id).all()
    for subscription in subscriptions:
        if subscription.renewal_date - datetime.now().date() <= timedelta(days=3):
            send_reminder(subscription)
            send_notification(subscription)

if __name__ == '__main__':
    scheduler.start()
    app.run(debug=True)
