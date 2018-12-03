import random
from nltk.corpus import movie_reviews
from textblob.classifiers import NaiveBayesClassifier
random.seed(1)
import simplejson as json
from flask import request
import flask
import sys
import nltk
nltk.download('punkt')
nltk.download('movie_reviews')
from nltk import word_tokenize,sent_tokenize
app = flask.Flask(__name__)
app.config["DEBUG"] = True

train = [
    ('I love this sandwich.', 'pos'),
    ('This is an amazing place!', 'pos'),
    ('I feel very good about these beers.', 'pos'),
    ('This is my best work.', 'pos'),
    ("What an awesome view", 'pos'),
    ('I do not like this restaurant', 'neg'),
    ('I am tired of this stuff.', 'neg'),
    ("I can't deal with this", 'neg'),
    ('He is my sworn enemy!', 'neg'),
    ('My boss is horrible.', 'neg')
]
test = [
    ('The beer was good.', 'pos'),
    ('I do not enjoy my job', 'neg'),
    ("I ain't feeling dandy today.", 'neg'),
    ("I feel amazing!", 'pos'),
    ('Gary is a friend of mine.', 'pos'),
    ("I can't believe I'm doing this.", 'neg')
]
cl = NaiveBayesClassifier(train)
reviews = [(list(movie_reviews.words(fileid)), category)
            for category in movie_reviews.categories()
            for fileid in movie_reviews.fileids(category)]
random.shuffle(reviews)
new_train, new_test = reviews[0:100], reviews[101:200]
cl.update(new_train)
@app.route('/', methods=['POST'])
def home():
    data = request.data
    dataDict = json.loads(data)
    return  cl.classify(dataDict['text'])

@app.route("/HEALTH")
def health():
   
    return "HEALTH"

app.run()



