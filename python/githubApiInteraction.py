#!/usr/bin/env python

import requests, re, sys, os

base = 'https://api.github.com/'
token = os.getenv('GHTOKEN')

if not token:
    message = """GH Token Unset
Your Github token needs to be set to interact with their
api and not hit the rate limit.
export GHTOKEN=<your github token>"""
    sys.exit(message)

class ApiError(Exception):
    def __init__(self, status):
        self.status = status
    def __str__(self):
        return "APIError: status={}".format(self.status)

def callGithubApi(api):
    resp = requests.get(base+api)
    if resp.status_code != 200:
        raise ApiError('Get /'+api+' {}'.format(resp.status_code))
    return (resp, api)

def userRepos(response):
    i=0
    for responseItem in response[0].json():
        print('{} {}'.format(responseItem['name'], responseItem['owner']['login']))
        i+=1
    print '\n{} has {} repositories'.format(re.findall(r'\/(.*?)\/', response[1])[0], str(i))

def repoCommits(response):
    i=0
    for item in response[0].json():
        author = item['author']
        if author:
            print '{}'.format(author['login'])
            i+=1
    print i
    return i

def getAllAuthorsForARepo(org, repo):
    total_count = 0
    page_num = 0
    while total_count < 13061:
        page_num+=1
        print "\nLooking at page: {}\nCurrent Total Count: {:,}".format(str(page_num), total_count)
        repoCommitsApi = 'repos/{}/{}/commits'.format(org, repo)
        queryString = '?page={}&per_page=100&access_token={}'.format(str(page_num), token)
        page_total = repoCommits(callGithubApi(repoCommitsApi+queryString))
        if page_total == 0:
            sys.exit("No authors found on this page.")
        else:
            total_count += page_total

# userOrgs = 'users/aln787/repos'
# userRepos(callGithubApi(userOrgs))
getAllAuthorsForARepo('fastlane', 'fastlane')
