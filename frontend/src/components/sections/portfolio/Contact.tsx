import { useState } from 'react'
import { motion } from 'framer-motion'
import { Send, Mail, Github, Linkedin, CheckCircle, Loader2 } from 'lucide-react'
import { api } from '../../../utils/apiConfig'

const socialLinks = [
  { icon: Github, href: 'https://github.com/jadenrazo', label: 'GitHub' },
  { icon: Linkedin, href: 'https://jadenrazo.dev/s/linkedin', label: 'LinkedIn' },
  { icon: Mail, href: 'mailto:contact@jadenrazo.dev', label: 'Email' },
]

export default function Contact() {
  const [isSubmitted, setIsSubmitted] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    subject: '',
    message: '',
    website: '',
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError(null)

    try {
      await api.post('/api/v1/contact', {
        name: formData.name,
        email: formData.email,
        subject: formData.subject,
        message: formData.message,
        website: formData.website,
      }, { skipAuth: true })

      setIsSubmitted(true)
      setFormData({ name: '', email: '', subject: '', message: '', website: '' })
      setTimeout(() => setIsSubmitted(false), 5000)
    } catch (err) {
      setError(
        err instanceof Error && err.message !== 'Request failed'
          ? err.message
          : 'Failed to send message. Please try again later.'
      )
    } finally {
      setIsLoading(false)
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setFormData((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }))
    if (error) setError(null)
  }

  return (
    <section id="contact" className="relative w-full py-12 sm:py-16 md:py-20">
      <div className="w-full flex flex-col items-center px-4 sm:px-6 lg:px-8">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6 }}
          className="text-center mb-6"
        >
          <h2 className="text-2xl sm:text-3xl md:text-4xl font-bold mb-2 text-text-primary">
            Let's <span className="gradient-text">Connect</span>
          </h2>
          <p className="text-text-secondary text-sm sm:text-base">
            Have a project in mind? I'd love to hear about it.
          </p>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6, delay: 0.1 }}
          className="glass-card p-4 sm:p-6 w-full max-w-xl"
        >
          {isSubmitted ? (
            <motion.div
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              className="flex flex-col items-center justify-center py-8 text-center"
            >
              <CheckCircle className="w-10 h-10 sm:w-12 sm:h-12 text-primary mb-3" />
              <h3 className="text-lg sm:text-xl font-bold mb-1 text-text-primary">Message Sent!</h3>
              <p className="text-text-secondary text-sm">
                Thanks for reaching out. I'll get back to you soon.
              </p>
            </motion.div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-3 sm:space-y-4">
              {error && (
                <div className="px-3 py-2 bg-red-500/10 border border-red-500/30 rounded-lg text-red-400 text-sm">
                  {error}
                </div>
              )}
              <div aria-hidden="true" style={{ position: 'absolute', left: '-9999px', opacity: 0, height: 0, overflow: 'hidden' }}>
                <input
                  type="text"
                  name="website"
                  value={formData.website}
                  onChange={handleChange}
                  tabIndex={-1}
                  autoComplete="off"
                />
              </div>
              <div className="grid sm:grid-cols-2 gap-3 sm:gap-4">
                <div>
                  <label className="sr-only" htmlFor="contact-name">Your name</label>
                  <input
                    type="text"
                    id="contact-name"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    required
                    disabled={isLoading}
                    className="w-full px-3 py-2.5 bg-surface border border-border rounded-lg focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-colors text-sm text-text-primary placeholder:text-text-muted disabled:opacity-50"
                    placeholder="Your name"
                  />
                </div>
                <div>
                  <label className="sr-only" htmlFor="contact-email">Your email</label>
                  <input
                    type="email"
                    id="contact-email"
                    name="email"
                    value={formData.email}
                    onChange={handleChange}
                    required
                    disabled={isLoading}
                    className="w-full px-3 py-2.5 bg-surface border border-border rounded-lg focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-colors text-sm text-text-primary placeholder:text-text-muted disabled:opacity-50"
                    placeholder="Your email"
                  />
                </div>
              </div>
              <div>
                <label className="sr-only" htmlFor="contact-subject">Subject</label>
                <input
                  type="text"
                  id="contact-subject"
                  name="subject"
                  value={formData.subject}
                  onChange={handleChange}
                  disabled={isLoading}
                  className="w-full px-3 py-2.5 bg-surface border border-border rounded-lg focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-colors text-sm text-text-primary placeholder:text-text-muted disabled:opacity-50"
                  placeholder="Subject (optional)"
                />
              </div>
              <div>
                <label className="sr-only" htmlFor="contact-message">Your message</label>
                <textarea
                  id="contact-message"
                  name="message"
                  value={formData.message}
                  onChange={handleChange}
                  required
                  disabled={isLoading}
                  rows={3}
                  className="w-full px-3 py-2.5 bg-surface border border-border rounded-lg focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-colors resize-none text-sm text-text-primary placeholder:text-text-muted disabled:opacity-50"
                  placeholder="Tell me about your project..."
                />
              </div>
              <button type="submit" className="btn-primary w-full" disabled={isLoading}>
                {isLoading ? (
                  <>
                    <Loader2 size={16} className="animate-spin" />
                    <span>Sending...</span>
                  </>
                ) : (
                  <>
                    <Send size={16} />
                    <span>Send Message</span>
                  </>
                )}
              </button>
            </form>
          )}
        </motion.div>

        <motion.div
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          viewport={{ once: true }}
          transition={{ delay: 0.4, duration: 0.5 }}
          className="flex items-center justify-center gap-4 mt-5"
        >
          {socialLinks.map((social) => (
            <a
              key={social.label}
              href={social.href}
              target="_blank"
              rel="noopener noreferrer"
              className="p-2.5 glass rounded-lg hover:bg-surface-hover transition-colors group"
              aria-label={social.label}
            >
              <social.icon className="w-5 h-5 text-text-secondary group-hover:text-primary transition-colors" />
            </a>
          ))}
        </motion.div>

        <motion.p
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          viewport={{ once: true }}
          transition={{ delay: 0.5, duration: 0.5 }}
          className="text-text-muted text-xs mt-5 text-center"
        >
          &copy; {new Date().getFullYear()} Jaden Razo. All rights reserved.
        </motion.p>
      </div>
    </section>
  )
}
