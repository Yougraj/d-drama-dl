import sys
import os
from flask import Flask, render_template, request, redirect, url_for
from urllib.parse import quote, unquote
import re
sys.path.insert(0, os.path.dirname(__file__))

# Import functions from main.py
from main import (
    fetch_dramas_from_page,
    get_random_drama,
    search_dramas,
    get_drama_episodes,
    get_episode_video_url,
    title_from_drama_link,
    BASE_URL,
)

app = Flask(__name__)

@app.route('/')
def index():
    """Homepage: Displays a selection of dramas from the main page."""
    print("Fetching dramas for homepage...")
    homepage_dramas = fetch_dramas_from_page(BASE_URL)
    # Limit to a reasonable number for display on the homepage
    display_dramas = homepage_dramas[:24] # Display first 24 dramas

    return render_template('index.html', dramas=display_dramas)

@app.route('/search', methods=['GET'])
def search():
    """Search page: Handles keyword and filter search, displays results with pagination."""
    keyword = request.args.get('keyword', '').strip()
    filter_type = request.args.get('filter', 'all').strip().lower()
    page = int(request.args.get('page', 1))
    
    search_results = []
    has_more_results = False

    if keyword:
        print(f"Web search for '{keyword}' (filter: {filter_type}, page: {page})...")
        results = search_dramas(keyword, filter_type, page)
        search_results.extend(results)
        
        # Simple pagination logic: if we got a full page of results (assuming 10 per page),
        # there might be more. This is a heuristic based on the AJAX endpoint's behavior.
        if len(results) == 10:
            has_more_results = True

    return render_template('search.html',
                           keyword=keyword,
                           filter_type=filter_type,
                           search_results=search_results,
                           current_page=page,
                           has_more_results=has_more_results)

@app.route('/drama/<path:drama_link_encoded>')
@app.route('/drama/<path:drama_link_encoded>/<path:image_url_encoded>')
def drama_detail(drama_link_encoded, image_url_encoded=None):
    """Drama detail page: Displays drama information and its episodes."""
    drama_link = unquote(drama_link_encoded)
    drama_image_url = unquote(image_url_encoded) if image_url_encoded else 'N/A'
    print(f"DEBUG: drama_link_encoded received: {drama_link_encoded}")
    print(f"DEBUG: image_url_encoded received: {image_url_encoded}")
    
    print(f"Fetching details and episodes for: {drama_link}")

    # Attempt to extract a title from the URL for display
    drama_title = title_from_drama_link(drama_link) or "Drama Details"
    latest_episode = request.args.get('latest_episode', 'N/A')

    episodes = get_drama_episodes(
        drama_link,
        drama_title=drama_title,
        latest_episode_label=latest_episode,
    )
    print(f"DEBUG: drama_image_url passed to template: {drama_image_url}")
    print(f"DEBUG: Number of episodes found: {len(episodes)}")
    
    return render_template('drama_detail.html',
                           drama_link=drama_link,
                           drama_title=drama_title,
                           drama_image_url=drama_image_url,
                           episodes=episodes)

@app.route('/watch')
def watch_episode():
    """Play a scraped episode video in the browser."""
    drama_link = unquote(request.args.get('drama_link', '')).strip()
    drama_title = unquote(request.args.get('drama_title', '')).strip()
    episode = request.args.get('episode', '1')

    if not drama_link:
        return redirect(url_for('search'))

    try:
        episode_num = int(episode)
    except ValueError:
        episode_num = 1

    if not drama_title:
        drama_title = title_from_drama_link(drama_link) or "Episode"

    episode_link = f"{drama_link.split('?')[0].rstrip('/')}/?episode={episode_num}"
    video_url = request.args.get('video_url', '').strip()
    if not video_url:
        video_url = get_episode_video_url(drama_title, episode_num)

    return render_template(
        'watch.html',
        drama_link=drama_link,
        drama_title=drama_title,
        episode_num=episode_num,
        episode_link=episode_link,
        video_url=video_url,
    )

if __name__ == '__main__':
    app.run(debug=True)