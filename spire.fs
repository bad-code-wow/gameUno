#version 100

precision mediump float;

// Input vertex attributes (from vertex shader)
varying vec2 fragTexCoord;
varying vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;
uniform vec2 points[16];

// NOTE: Add your custom variables here

const vec2 size = vec2(800, 450);   // Framebuffer size




void main()
{



    // Texel color fetching from texture sampler
    vec4 source = texture2D(texture0, fragTexCoord);

    vec3 col = vec3(fragTexCoord.x - points.x,fragTexCoord.y - points.y,0.0);
    

    gl_FragColor = vec4(col,1);

}